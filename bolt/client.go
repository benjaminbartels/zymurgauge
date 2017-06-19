package bolt

import (
	"encoding/binary"
	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client allows communication with the BoltBD datastore
type Client struct {
	Path                string
	logger              *logrus.Logger
	db                  *bolt.DB
	chamberService      ChamberService
	fermentationService FermentationService
	beerService         BeerService
}

// NewClient creates a new client using the given path to the BoltDB datastore
func NewClient(path string, logger *logrus.Logger) *Client {
	return &Client{
		Path:   path,
		logger: logger,
	}
}

// Open opens the connection to the BoltDB datastore and initializes the buckets if necessary
func (c *Client) Open() error {
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrapf(err, "Could not open %s", c.Path)
	}
	c.db = db

	tx, err := c.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "Could not begin transaction")
	}

	defer func() { _ = tx.Rollback() }()

	if _, err := tx.CreateBucketIfNotExists([]byte("Chambers")); err != nil {
		return errors.Wrap(err, "Could not create Chambers bucket")
	}

	if _, err := tx.CreateBucketIfNotExists([]byte("Fermentations")); err != nil {
		return errors.Wrap(err, "Could not create Fermentations bucket")
	}

	if _, err := tx.CreateBucketIfNotExists([]byte("Beers")); err != nil {
		return errors.Wrap(err, "Could not create Beer bucket")
	}

	c.chamberService.db = db
	c.fermentationService.db = db
	c.beerService.db = db
	c.chamberService.logger = c.logger
	c.fermentationService.logger = c.logger
	c.beerService.logger = c.logger

	return tx.Commit()
}

// Close closes the connection to the BoltDB datastore
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ChamberService returns the service used to manage Chambers
func (c *Client) ChamberService() zymurgauge.ChamberService {
	return &c.chamberService
}

// FermentationService returns the service used to manage Fermentations
func (c *Client) FermentationService() zymurgauge.FermentationService {
	return &c.fermentationService
}

// BeerService returns the service used to manage Beers
func (c *Client) BeerService() zymurgauge.BeerService {
	return &c.beerService
}

// itob encodes unsigned 64-bit integer to byte slices to improve performance when used as BoltDB keys
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
