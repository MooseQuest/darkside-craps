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
	verifs   *mongo.Collection // pending email-verification codes, keyed by email
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
		verifs:   db.Collection("email_verifications"),
	}
	// Best-effort indexes.
	_, _ = s.users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true),
	})
	_, _ = s.sessions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}, {Key: "created_at", Value: -1}},
	})
	// Stale, unfinished ceremonies self-expire so anonymous begin calls can't
	// accumulate documents indefinitely.
	_, _ = s.waSess.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(600),
	})
	// Verification codes expire exactly at their expires_at time (TTL 0).
	_, _ = s.verifs.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0),
	})
	return s, nil
}

func newHandle() ([]byte, error) {
	h := make([]byte, 32)
	_, err := rand.Read(h)
	return h, err
}

// ---- User model ------------------------------------------------------------

type userDoc struct {
	Email     string    `bson:"email"`
	Handle    []byte    `bson:"handle"` // WebAuthn user id
	CreatedAt time.Time `bson:"created_at"`
}

type credDoc struct {
	Email string              `bson:"email"`
	Cred  webauthn.Credential `bson:"cred"`
}

// User implements webauthn.User.
type User struct {
	email  string
	handle []byte
	creds  []webauthn.Credential
}

func (u *User) WebAuthnID() []byte                         { return u.handle }
func (u *User) WebAuthnName() string                       { return u.email }
func (u *User) WebAuthnDisplayName() string                { return u.email }
func (u *User) WebAuthnCredentials() []webauthn.Credential { return u.creds }
func (u *User) WebAuthnIcon() string                       { return "" }

// GetUser returns an existing user or errNotFound. It never creates a record —
// user documents are only persisted once a passkey registration finishes.
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

// EnsureUser persists the user with the given handle if it does not yet exist,
// and returns the (existing or new) user. Called only at registration finish.
func (s *Store) EnsureUser(ctx context.Context, email string, handle []byte) (*User, error) {
	var doc userDoc
	err := s.users.FindOneAndUpdate(
		ctx,
		bson.M{"email": email},
		bson.M{"$setOnInsert": userDoc{Email: email, Handle: handle, CreatedAt: time.Now().UTC()}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
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
	Handle    []byte               `bson:"handle"` // user handle chosen at begin, reused at finish
	CreatedAt time.Time            `bson:"created_at"`
}

func (s *Store) SaveWASession(ctx context.Context, email string, data *webauthn.SessionData, handle []byte) error {
	_, err := s.waSess.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": waSessionDoc{Email: email, Data: *data, Handle: handle, CreatedAt: time.Now().UTC()}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *Store) TakeWASession(ctx context.Context, email string) (*webauthn.SessionData, []byte, error) {
	var doc waSessionDoc
	err := s.waSess.FindOneAndDelete(ctx, bson.M{"email": email}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil, errNotFound
	} else if err != nil {
		return nil, nil, err
	}
	return &doc.Data, doc.Handle, nil
}

// ---- Email verification codes ---------------------------------------------

type verifyDoc struct {
	Email     string    `bson:"email"`
	CodeHash  string    `bson:"code_hash"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

// SaveVerification stores (or replaces) the pending code hash for an email.
func (s *Store) SaveVerification(ctx context.Context, email, codeHash string, ttl time.Duration) error {
	now := time.Now().UTC()
	_, err := s.verifs.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": verifyDoc{Email: email, CodeHash: codeHash, ExpiresAt: now.Add(ttl), CreatedAt: now}},
		options.Update().SetUpsert(true),
	)
	return err
}

// ConsumeVerification atomically checks a code hash for an email and deletes it
// (single-use). Returns true only for a matching, unexpired code.
func (s *Store) ConsumeVerification(ctx context.Context, email, codeHash string) (bool, error) {
	var doc verifyDoc
	err := s.verifs.FindOneAndDelete(ctx, bson.M{
		"email":      email,
		"code_hash":  codeHash,
		"expires_at": bson.M{"$gt": time.Now().UTC()},
	}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ---- Game history ----------------------------------------------------------

// gameSummary is the whitelisted, strictly-typed shape we persist. Only numeric
// fields are accepted, so a client cannot store arbitrary/oversized JSON or
// smuggle markup into the history view.
type gameSummary struct {
	InitialBankroll   float64 `json:"initial_bankroll" bson:"initial_bankroll"`
	FinalBankroll     float64 `json:"final_bankroll" bson:"final_bankroll"`
	RollCount         int     `json:"roll_count" bson:"roll_count"`
	PointsEstablished int     `json:"points_established" bson:"points_established"`
	PointsHit         int     `json:"points_hit" bson:"points_hit"`
	PointsLost        int     `json:"points_lost" bson:"points_lost"`
	Wins              float64 `json:"wins" bson:"wins"`
	Losses            float64 `json:"losses" bson:"losses"`
	Strikes           int     `json:"strikes" bson:"strikes"`
	TotalWagered      float64 `json:"total_wagered" bson:"total_wagered"`
	Net               float64 `json:"net" bson:"net"`
}

type gameSession struct {
	Email     string      `bson:"email" json:"-"`
	Summary   gameSummary `bson:"summary" json:"summary"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
}

func (s *Store) SaveGameSession(ctx context.Context, email string, summary gameSummary) error {
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
