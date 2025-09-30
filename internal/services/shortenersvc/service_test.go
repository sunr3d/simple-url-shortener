package shortenersvc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sunr3d/simple-url-shortener/mocks"
	"github.com/sunr3d/simple-url-shortener/models"
)

const (
	testLink    = "https://example.com/x"
	baseTestURL = "http://TESTshrt.ly"
)

// ShortenLink Tests.
func Test_ShortenLink_OK(t *testing.T) {
	db := mocks.NewDatabase(t)

	db.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(l *models.Link) bool {
			return l != nil && len(l.Code) == codeLen && l.Original == testLink
		})).
		Return(nil)

	svc := New(db, baseTestURL)

	code, shortURL, err := svc.ShortenLink(context.Background(), testLink)
	require.NoError(t, err)
	require.Len(t, code, codeLen)
	require.Equal(t, baseTestURL+"/s/"+code, shortURL)
}

func Test_ShortenLink_InvalidScheme(t *testing.T) {
	db := mocks.NewDatabase(t)
	svc := New(db, baseTestURL)

	_, _, err := svc.ShortenLink(context.Background(), "ftp://example.com/file")
	require.Error(t, err)

	db.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func Test_ShortenLink_RepoErr(t *testing.T) {
	db := mocks.NewDatabase(t)

	db.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*models.Link")).
		Return(errors.New("repo error"))

	svc := New(db, baseTestURL)

	_, _, err := svc.ShortenLink(context.Background(), testLink)
	require.Error(t, err)
}

// FollowLink Tests.
func Test_FollowLink_OK(t *testing.T) {
	db := mocks.NewDatabase(t)

	db.EXPECT().
		GetLink(mock.Anything, "abcde").
		Return(&models.Link{Code: "abcde", Original: testLink}, nil)

	svc := New(db, baseTestURL)

	orig, err := svc.FollowLink(context.Background(), "abcde")
	require.NoError(t, err)
	require.Equal(t, testLink, orig)
}

func Test_FollowLink_EmptyCode(t *testing.T) {
	db := mocks.NewDatabase(t)
	svc := New(db, baseTestURL)

	_, err := svc.FollowLink(context.Background(), "")
	require.Error(t, err)
}

func Test_FollowLink_NotFound(t *testing.T) {
	db := mocks.NewDatabase(t)

	db.EXPECT().
		GetLink(mock.Anything, "1a2b3c").
		Return(nil, nil)

	svc := New(db, baseTestURL)

	_, err := svc.FollowLink(context.Background(), "1a2b3c")
	require.Error(t, err)
}
