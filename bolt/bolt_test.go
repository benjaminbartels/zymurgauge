package bolt_test

import (
	"io/ioutil"
	"os"

	"github.com/orangesword/zymurgauge/bolt"
	"github.com/sirupsen/logrus"
)

// TestClient is a wrapper around the bolt.Client.
type TestClient struct {
	*bolt.Client
}

// MustOpenClient returns a new TestClient with an open underlying bolt.Client and creates a temp datatore
func MustOpenClient() *TestClient {

	l := logrus.New()

	p := "zymurgauge-bolt-client-"
	f, err := ioutil.TempFile("", p)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	c := &TestClient{
		Client: bolt.NewClient(p, l),
	}
	c.Path = f.Name()

	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close removes the temp bolt datastore and closes the underlying bolt.Client
func (c *TestClient) Close() error {
	defer func() { _ = os.Remove(c.Path) }()
	return c.Client.Close()
}
