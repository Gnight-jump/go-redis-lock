package alock

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func LockTransaction(key, reqID string, f func() error) (bool, error) {
	lock := NewRLock(key, reqID, 10*time.Second)

	result := lock.TryLock()
	if !result {
		log.Errorf("lock exist, key=%v, reqID=%v", key, reqID)
		return false, nil
	}

	stop := make(chan bool, 1)
	defer func() {
		stop <- true
		close(stop)
		err := lock.UnLock()
		if err != nil {
			log.Errorf("unlock failed, key=%v, reqID=%v", key, reqID)
			return
		}
		log.Infof("unlock succeed, key=%v, reqID=%v", key, reqID)
	}()

	// automatic renewal of additional maintenance lock
	go renewlRLock(lock, stop, key, reqID)
	return true, f()
}

func renewlRLock(lock *RLock, stop chan bool, key, reqID string) {
	ticker := time.NewTicker(5 * time.Second)
	renewTime := 0
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			log.Errorf("renew RLock happened, key=%v, reqID=%v, err=%v", key, reqID, err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			err := lock.RenewLock()
			if err != nil {
				log.Errorf("renew RLock failed, key=%v, reqID=%v, err=%v", key, reqID, err)
				return
			}
			renewTime++
			log.Infof("renew RLock succeed, key=%v, reqID=%v, renewTime=%v", key, reqID, renewTime)
		case <-stop:
			log.Infof("renew RLock over, key=%v, reqID=%v", key, reqID)
			return
		}
	}
}
