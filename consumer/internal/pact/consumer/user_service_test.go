//go:build integration

package consumer

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"net/url"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pact-foundation/pact-workshop-go/consumer/internal/client"
	"github.com/pact-foundation/pact-workshop-go/consumer/internal/model"
	"github.com/stretchr/testify/assert"
)

var Like = matchers.Like
var EachLike = matchers.EachLike
var Term = matchers.Term
var Regex = matchers.Regex
var HexValue = matchers.HexValue
var Identifier = matchers.Identifier
var IPAddress = matchers.IPAddress
var IPv6Address = matchers.IPv6Address
var Timestamp = matchers.Timestamp
var Date = matchers.Date
var Time = matchers.Time
var UUID = matchers.UUID
var ArrayMinLike = matchers.ArrayMinLike

type S = matchers.S
type Map = matchers.MapMatcher

var u *url.URL
var apiClient *client.Client

func TestMain(m *testing.M) {
	log.SetLogLevel("INFO")
	os.Exit(m.Run())
}

func TestUserServicePact_GetUser(t *testing.T) {
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: os.Getenv("CONSUMER_NAME"),
		Provider: os.Getenv("PROVIDER_NAME"),
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
	})
	assert.NoError(t, err)

	t.Run("the user exists", func(t *testing.T) {
		id := 10

		err = mockProvider.
			AddInteraction().
			Given("User sally exists").
			UponReceiving("A request to login with user 'sally'").
			WithRequestPathMatcher("GET", Regex("/user/"+strconv.Itoa(id), "/user/[0-9]+"), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", Like("Bearer 2019-01-01"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.BodyMatch(model.User{}).
					Header("Content-Type", Term("application/json", `application\/json`)).
					Header("X-Api-Correlation-Id", Like("100"))
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Act: test our API client behaves correctly

				// Get the Pact mock server URL
				u, _ = url.Parse("http://" + config.Host + ":" + strconv.Itoa(config.Port))

				// Initialise the API client and point it at the Pact mock server
				apiClient = &client.Client{
					BaseURL: u,
				}

				// // Execute the API client
				user, err := apiClient.WithToken("2019-01-01").GetUser(id)

				// // Assert basic fact
				if user.ID != id {
					return fmt.Errorf("wanted user with ID %d but got %d", id, user.ID)
				}

				return err
			})

		assert.NoError(t, err)

	})

	t.Run("the user does not exist", func(t *testing.T) {
		id := 10

		err = mockProvider.
			AddInteraction().
			Given("User sally does not exist").
			UponReceiving("A request to login with user 'sally'").
			WithRequestPathMatcher("GET", Regex("/user/"+strconv.Itoa(id), "/user/[0-9]+"), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", Like("Bearer 2019-01-01"))
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", Term("application/json", `application\/json`)).
					Header("X-Api-Correlation-Id", Like("100"))
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Act: test our API client behaves correctly

				// Get the Pact mock server URL
				u, _ = url.Parse("http://" + config.Host + ":" + strconv.Itoa(config.Port))

				// Initialise the API client and point it at the Pact mock server
				apiClient = &client.Client{
					BaseURL: u,
				}

				// // Execute the API client
				_, err := apiClient.WithToken("2019-01-01").GetUser(id)
				assert.Equal(t, client.ErrNotFound, err)
				return nil
			})
		assert.NoError(t, err)

	})
	t.Run("not authenticated", func(t *testing.T) {
		id := 10

		err = mockProvider.
			AddInteraction().
			Given("User is not authenticated").
			UponReceiving("A request to login with user 'sally'").
			WithRequestPathMatcher("GET", Regex("/user/"+strconv.Itoa(id), "/user/[0-9]+")).
			WillRespondWith(401, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", Term("application/json", `application\/json`)).
					Header("X-Api-Correlation-Id", Like("100"))
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Act: test our API client behaves correctly

				// Get the Pact mock server URL
				u, _ = url.Parse("http://" + config.Host + ":" + strconv.Itoa(config.Port))

				// Initialise the API client and point it at the Pact mock server
				apiClient = &client.Client{
					BaseURL: u,
				}

				// // Execute the API client
				_, err := apiClient.WithToken("").GetUser(id)
				assert.Equal(t, client.ErrUnauthorized, err)
				return nil
			})
		assert.NoError(t, err)
	})

}

func TestUserServicePact_Health(t *testing.T) {
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: os.Getenv("CONSUMER_NAME"),
		Provider: os.Getenv("PROVIDER_NAME"),
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
	})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		Given("The service is alive").
		UponReceiving("A request to check the health of the service").
		WithRequestPathMatcher("GET", Regex("/healthasdf", "/healthasdf")).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", Term("text/plain", `text\/plain`))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			// Act: test our API client behaves correctly

			// Get the Pact mock server URL
			u, _ = url.Parse("http://" + config.Host + ":" + strconv.Itoa(config.Port))

			// Initialise the API client and point it at the Pact mock server
			apiClient = &client.Client{
				BaseURL: u,
			}

			// // Execute the API client
			err := apiClient.Health()
			assert.NoError(t, err)
			return nil
		})
	assert.NoError(t, err)
}
