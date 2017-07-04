package bolt

import (
	"encoding/json"
	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// check implementation at compile time
var _ zymurgauge.ChamberService = &ChamberService{}

// ChamberService represents a service for managing controllers
type ChamberService struct {
	db      *bolt.DB
	logger  *logrus.Logger
	clients map[string]chan zymurgauge.Chamber
}

func (s *ChamberService) GetAll() ([]zymurgauge.Chamber, error) {
	var chambers []zymurgauge.Chamber

	err := s.db.View(func(tx *bolt.Tx) error {

		if err := tx.Bucket([]byte("Chambers")).ForEach(func(k, v []byte) error {

			var c zymurgauge.Chamber
			if err := json.Unmarshal(v, &c); err != nil {
				return errors.Wrap(err, "Could not unmarshal Chambers")
			}

			chambers = append(chambers, c)

			return nil
		}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return chambers, nil
}

// Get returns a Chamber by its MAC address
func (s *ChamberService) Get(mac string) (*zymurgauge.Chamber, error) {
	var c zymurgauge.Chamber

	err := s.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket([]byte("Chambers")).Get([]byte(mac)); v == nil {
			return zymurgauge.ErrNotFound
		} else if err := json.Unmarshal(v, &c); err != nil {
			return errors.Wrapf(err, "Could not unmarshal Chamber %s", mac)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Save creates or updates a Chamber
func (s *ChamberService) Save(c *zymurgauge.Chamber) error {
	err := s.db.Update(func(tx *bolt.Tx) error {

		bu := tx.Bucket([]byte("Chambers"))
		c.ModTime = time.Now()

		if v, err := json.Marshal(c); err != nil {
			return errors.Wrapf(err, "Could not marshal Chamber %s", c.MacAddress)
		} else if err := bu.Put([]byte(c.MacAddress), v); err != nil {
			return errors.Wrapf(err, "Could not put Chamber %s", c.MacAddress)
		}

		go s.send(c)

		return nil
	})
	return err
}

// Subscribe registers the caller to receives updates to the given controller on the given channel
func (s *ChamberService) Subscribe(mac string, ch chan zymurgauge.Chamber) error {
	c, err := s.Get(mac)
	if err != nil {
		return errors.Wrapf(err, "Could not get Chamber %d", mac)
	}

	if s.clients == nil {
		s.clients = make(map[string]chan zymurgauge.Chamber)
	}

	s.clients[c.MacAddress] = ch

	return nil
}

// Unsubscribe unregisters the caller to receives updates to the given controller
func (s *ChamberService) Unsubscribe(mac string) {
	if s.clients != nil {
		ch, ok := s.clients[mac]
		if ok {
			close(ch)
			delete(s.clients, mac)
		}
	}
}

// send sends the Chamber to the listening client's channel
func (s *ChamberService) send(c *zymurgauge.Chamber) {
	ch, ok := s.clients[c.MacAddress]
	if ok {
		s.logger.Debugf("Sending update for Chamber %s", c.MacAddress)
		ch <- *c
	} else {
		s.logger.Debugf("No listeners for Chamber %s", c.MacAddress)
	}
}
