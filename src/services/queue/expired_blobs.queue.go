package services

import (
	"blob/src/database"
	"blob/src/functions"
	"blob/src/models"
	"blob/src/services"
	"context"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

const TypeDeleteExpiredBlobs = "blob:delete_expired"

func handleDeleteExpiredBlobs(ctx context.Context, t *asynq.Task) error {
	return removeExpiredBlobs()
}

func StartQueueWorker() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeDeleteExpiredBlobs, handleDeleteExpiredBlobs)
	go func() {
		if services.AsynqServer == nil {
			functions.Error("[QUEUE ERROR] AsynqServer is nil! Did you call InitAsynq() first?")
			return
		}
		if err := services.AsynqServer.Run(mux); err != nil {
			functions.Error("[QUEUE ERROR] Worker failed: %v", err)
		}
	}()
}

func StartCleanupScheduler() {

	interval := 24 * time.Hour
	if v := os.Getenv("BLOB_EXPIRED_CLEANUP_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		} else {
			functions.Warn("[QUEUE] Invalid BLOB_EXPIRED_CLEANUP_INTERVAL '%s', using default 24h", v)
		}
	}
	functions.Info("[QUEUE] Expired blob scheduler started (interval: %v)", interval)
	go func() {
		for {
			if services.AsynqClient == nil {
				functions.Error("[QUEUE ERROR] AsynqClient is nil! Did you call InitAsynq() first?")
				return
			}
			task := asynq.NewTask(TypeDeleteExpiredBlobs, nil)
			_, err := services.AsynqClient.Enqueue(task)
			if err != nil {
				functions.Error("[QUEUE ERROR] Failed to enqueue blob expired task: %v", err)
			}
			time.Sleep(interval)
		}
	}()
}

func removeExpiredBlobs() error {
	var expired []models.Blob
	err := database.DB.Where("expires_at IS NOT NULL AND expires_at <= ?", time.Now()).Find(&expired).Error
	if err != nil {
		functions.Error("[QUEUE] Failed to query expired blobs: %v", err)
		return err
	}
	for _, blob := range expired {
		if blob.Path != "" {
			functions.Info("[QUEUE] Attempting to remove file: %s", blob.Path)
			if stat, statErr := os.Stat(blob.Path); statErr != nil {
				functions.Warn("[QUEUE] File stat error for %s: %v", blob.Path, statErr)
			} else {
				functions.Info("[QUEUE] File exists. Size: %d, Mode: %v", stat.Size(), stat.Mode())
			}
			err := os.Remove(blob.Path)
			if err != nil {
				functions.Warn("[QUEUE] os.Remove error for %s: %v", blob.Path, err)
			} else {
				functions.Info("[QUEUE] File removed: %s", blob.Path)
			}
		}
		delErr := database.DB.Delete(&blob).Error
		if delErr != nil {
			functions.Warn("[QUEUE] Failed to remove blob from DB: %v", delErr)
		} else {
			functions.Info("[QUEUE] Blob removed from DB: %s", blob.ID.String())
		}
	}
	functions.Info("[QUEUE] Expired blobs expired finished. Total: %d", len(expired))
	return nil
}
