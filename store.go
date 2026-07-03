package main

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errNotFound = errors.New("not found")

// Store wraps the MongoDB collections the app uses.
type Store struct {
	client   *mongo.Client
	users    *mongo.Collection
	creds    *mongo.Collection
	waSess   *mongo.Collection // in-flight WebAuthn ceremony data, keyed by email
	sessions *mongo.Collection // saved game summaries
}

func NewStore(ctx context.Context, cfg Config) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := client.Database(cfg.DBName)
	s := &Store{
		client:   client,
		users:    db.Collection("users"),
		creds:    db.Collection("credentials"),
		waSess:   db.Collection("wa_sessions"),
		sessions: db.Collection("game_sessions"),
	}
	// Best-effort indexes.
	_, _ = s.users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true),
	})
	_, _ = s.sessions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}, {Key: "created_at", Value: -1}},
	})
	return s, nil
}

// ---- User model ------------------------------------------------------------

type userDoc struct {
	Email      string    `bson:"email"`
	Handle     []byte    `bson:"handle"` // WebAuthn user id
	CreatedAt  time.Time `bson:"created_at"`
}

type credDoc struct {
	Email string              `bson:"email"`
	Cred  webauthn.Credential `bson:"cred"`
}

// User implements webauthn.User.
type User struct {
	email string
	handle []byte
	creds  []webauthn.Credential
}

func (u *User) WebAuthnID() []byte                         { return u.handle }
func (u *User) WebAuthnName() string                       { return u.email }
func (u *User) WebAuthnDisplayName() string                { return u.email }
func (u *User) WebAuthnCredentials() []webauthn.Credential { return u.creds }
func (u *User) WebAuthnIcon() string                       { return "" }

// GetOrCreateUser returns the user, creating a bare record (no creds) if new.
func (s *Store) GetOrCreateUser(ctx context.Context, email string) (*User, error) {
	var doc userDoc
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		handle := make([]byte, 32)
		if _, err := rand.Read(handle); err != nil {
			return nil, err
		}
		doc = userDoc{Email: email, Handle: handle, CreatedAt: time.Now().UTC()}
		if _, err := s.users.InsertOne(ctx, doc); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	creds, err := s.credentialsFor(ctx, email)
	if err != nil {
		return nil, err
	}
	return &User{email: email, handle: doc.Handle, creds: creds}, nil
}

// GetUser returns an existing user or errNotFound.
func (s *Store) GetUser(ctx context.Context, email string) (*User, error) {
	var doc userDoc
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, errNotFound
	} else if err != nil {
		return nil, err
	}
	creds, err := s.credentialsFor(ctx, email)
	if err != nil {
		return nil, err
	}
	return &User{email: email, handle: doc.Handle, creds: creds}, nil
}

func (s *Store) credentialsFor(ctx context.Context, email string) ([]webauthn.Credential, error) {
	cur, err := s.creds.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []webauthn.Credential
	for cur.Next(ctx) {
		var cd credDoc
		if err := cur.Decode(&cd); err != nil {
			return nil, err
		}
		out = append(out, cd.Cred)
	}
	return out, cur.Err()
}

func (s *Store) AddCredential(ctx context.Context, email string, c *webauthn.Credential) error {
	_, err := s.creds.InsertOne(ctx, credDoc{Email: email, Cred: *c})
	return err
}

// UpdateCredential bumps the stored sign counter after a successful login.
func (s *Store) UpdateCredential(ctx context.Context, email string, c *webauthn.Credential) error {
	_, err := s.creds.UpdateOne(ctx,
		bson.M{"email": email, "cred.id": c.ID},
		bson.M{"$set": bson.M{"cred": *c}},
	)
	return err
}

// ---- WebAuthn ceremony session data ---------------------------------------

type waSessionDoc struct {
	Email     string               `bson:"email"`
	Data      webauthn.SessionData `bson:"data"`
	CreatedAt time.Time            `bson:"created_at"`
}

func (s *Store) SaveWASession(ctx context.Context, email string, data *webauthn.SessionData) error {
	_, err := s.waSess.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": waSessionDoc{Email: email, Data: *data, CreatedAt: time.Now().UTC()}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *Store) TakeWASession(ctx context.Context, email string) (*webauthn.SessionData, error) {
	var doc waSessionDoc
	err := s.waSess.FindOneAndDelete(ctx, bson.M{"email": email}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, errNotFound
	} else if err != nil {
		return nil, err
	}
	return &doc.Data, nil
}

// ---- Game history ----------------------------------------------------------

type gameSession struct {
	Email     string         `bson:"email" json:"-"`
	Summary   map[string]any `bson:"summary" json:"summary"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
}

func (s *Store) SaveGameSession(ctx context.Context, email string, summary map[string]any) error {
	_, err := s.sessions.InsertOne(ctx, gameSession{
		Email: email, Summary: summary, CreatedAt: time.Now().UTC(),
	})
	return err
}

func (s *Store) History(ctx context.Context, email string, limit int64) ([]gameSession, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cur, err := s.sessions.Find(ctx, bson.M{"email": email}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []gameSession
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
