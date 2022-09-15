package commands

import (
	"time"

	"github.com/jimeh/go-midjourney"
	"github.com/spf13/cobra"
)

func NewMidjourneyRecentJobs(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "recent-jobs",
		Aliases: []string{"jobs", "recent", "rj", "j", "r"},
		Short:   "List recent jobs",
		RunE:    midjourneyRecentJobsRunE(mc),
	}

	cmd.Flags().StringP("format", "f", "", "output format (yaml or json)")
	cmd.Flags().IntP("amount", "a", 50, "amount of jobs to list")
	cmd.Flags().StringP("type", "t", "", "type of jobs to list")
	cmd.Flags().StringP("order", "o", "new", "either \"new\" or \"oldest\"")
	cmd.Flags().StringP("user-id", "u", "", "user ID to list jobs for")
	cmd.Flags().IntP("page", "p", 0, "page to fetch")
	cmd.Flags().StringP("prompt", "s", "", "prompt text to search for")
	cmd.Flags().Bool("dedupe", true, "dedupe results")

	return cmd, nil
}

func midjourneyRecentJobsRunE(mc *midjourney.Client) runEFunc {
	return func(cmd *cobra.Command, _ []string) error {
		fs := cmd.Flags()
		q := &midjourney.RecentJobsQuery{}

		if v, err := fs.GetInt("amount"); err == nil && v > 0 {
			q.Amount = v
		}
		if v, err := fs.GetString("type"); err == nil && v != "" {
			q.JobType = midjourney.JobType(v)
		}
		if v, err := fs.GetString("order"); err == nil && v != "" {
			q.OrderBy = midjourney.Order(v)
		}
		if v, err := fs.GetString("user-id"); err == nil && v != "" {
			q.UserID = v
		}
		if v, err := fs.GetInt("page"); err == nil && v != 0 {
			q.Page = v
		}
		if v, err := fs.GetString("prompt"); err == nil && v != "" {
			q.Prompt = v
		}
		if v, err := fs.GetBool("dedupe"); err == nil {
			q.Dedupe = v
		}

		rj, err := mc.RecentJobs(cmd.Context(), q)
		if err != nil {
			return err
		}

		r := []*MidjourneyJob{}
		for _, j := range rj.Jobs {
			r = append(r, &MidjourneyJob{
				ID:             j.ID,
				Status:         string(j.CurrentStatus),
				Type:           string(j.Type),
				EnqueueTime:    j.EnqueueTime.Time,
				Prompt:         j.Prompt,
				ImagePaths:     j.ImagePaths,
				IsPublished:    j.IsPublished,
				UserID:         j.UserID,
				Username:       j.Username,
				FullCommand:    j.FullCommand,
				ReferenceJobID: j.ReferenceJobID,
			})
		}
		format := flagString(cmd, "format")

		return render(cmd.OutOrStdout(), format, r)
	}
}

type MidjourneyJob struct {
	ID             string    `json:"id,omitempty" yaml:"id,omitempty"`
	Status         string    `json:"current_status,omitempty" yaml:"current_status,omitempty"`
	Type           string    `json:"type,omitempty" yaml:"type,omitempty"`
	EnqueueTime    time.Time `json:"enqueue_time,omitempty" yaml:"enqueue_time,omitempty"`
	Prompt         string    `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	ImagePaths     []string  `json:"image_paths,omitempty" yaml:"image_paths,omitempty"`
	IsPublished    bool      `json:"is_published,omitempty" yaml:"is_published,omitempty"`
	UserID         string    `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Username       string    `json:"username,omitempty" yaml:"username,omitempty"`
	FullCommand    string    `json:"full_command,omitempty" yaml:"full_command,omitempty"`
	ReferenceJobID string    `json:"reference_job_id,omitempty" yaml:"reference_job_id,omitempty"`
}
