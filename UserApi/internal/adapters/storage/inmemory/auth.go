package inmemory

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	"context"
	"errors"

	"sync"
)

type SessionKey struct {
	UserID   string
	DeviceID string
}

type SessionRepository struct {
	mu sync.RWMutex

	sessions map[SessionKey]*authCore.SessionModel
}

func NewSessionStorage() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[SessionKey]*authCore.SessionModel),
	}
}

func (s *SessionRepository) Create(ctx context.Context,
	SessionModel authCore.SessionModel,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := SessionKey{
		UserID:   SessionModel.UserID,
		DeviceID: SessionModel.DeviceID,
	}
	s.sessions[key] = &SessionModel
	return nil
}

func (s *SessionRepository) GetByUserIdAndDeviceId(ctx context.Context,
	userID string,
	deviceID string,
) (
	authCore.SessionModel,
	error,
) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := SessionKey{
		UserID:   userID,
		DeviceID: deviceID,
	}

	sessionModel, ok := s.sessions[key]
	if !ok {
		return authCore.SessionModel{}, errors.New("session not found")
	}

	return *sessionModel, nil
}

func (s *SessionRepository) Update(ctx context.Context,
	SessionModel authCore.SessionModel,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := SessionKey{
		UserID:   SessionModel.UserID,
		DeviceID: SessionModel.DeviceID,
	}

	sessionModel, ok := s.sessions[key]
	if !ok {
		return errors.New("session not found")
	}

	*sessionModel = SessionModel

	return nil
}

func (s *SessionRepository) DeleteByUserAndDevice(ctx context.Context,
	userID,
	deviceID string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := SessionKey{
		UserID:   userID,
		DeviceID: deviceID,
	}

	_, ok := s.sessions[key]
	if !ok {
		return errors.New("session not found")
	}

	delete(s.sessions, key)

	return nil
}
