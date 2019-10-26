package jobs

/*


// UpdateWorkspaceJob updates the workspace metadata periodically
func UpdateWorkspaceJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	//auth := backend.GetAuthToken(ctx, id)

	logger.Info(ctx, "jobs.update.workspace", "workspace=%s", id)
}

// UpdateUsersJob updates the list of users of a workspace
func UpdateUsersJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	cursor := c.Query("cursor")
	auth := backend.GetAuthToken(ctx, id)

	logger.Info(ctx, "jobs.update.users", "workspace=%s", id)

	// update the list of users
	users, err := slack.UsersList(ctx, auth, cursor)

	if err == nil {

		if users.OK {

			logger.Info(ctx, "jobs.update.users", "users=%d", len(users.Members))
			metrics.Count(ctx, "jobs.update.users", id, len(users.Members))

			for i := range users.Members {
				err = backend.UpdateUser(ctx, users.Members[i].ID, users.Members[i].TeamID, users.Members[i].Name, users.Members[i].RealName, users.Members[i].Profile.FirstName, users.Members[i].Profile.LastName, users.Members[i].Profile.Email, users.Members[i].Deleted, users.Members[i].IsBot)

				if err != nil {
					logger.Error(ctx, "jobs.update.users", err.Error())
				} else {
					logger.Info(ctx, "jobs.update.users", "user=%s", users.Members[i].ID)
				}
			}

			nextCursor := users.ResponseMetadata["next_cursor"]
			if nextCursor != "" {
				// there is more data, schedule its retrieval
				job.ScheduleJob(ctx, backend.BackgroundWorkQueue, fmt.Sprintf(types.JobsBaseURL+"/users?id=%v&cursor=%v", id, nextCursor))

				logger.Info(ctx, "jobs.update.users", "next=%s", nextCursor)
			}
		} else {
			// Slack API returned an error
			logger.Critical(ctx, "jobs.update.users", "status=%s", users.Error)
		}

	} else {
		logger.Error(ctx, "jobs.update.users", err.Error())
	}
}

// UpdateChannelsJob updates the workspace metadata periodically
func UpdateChannelsJob(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	id := c.Query("id")
	cursor := c.Query("cursor")
	auth := backend.GetAuthToken(ctx, id)

	logger.Info(ctx, "jobs.update.channels", "workspace=%s", id)

	// update the list of channels
	channels, err := slack.ChannelsList(ctx, auth, cursor)
	if err == nil {

		if channels.OK {

			logger.Info(ctx, "jobs.update.channels", "channels=%d", len(channels.Channels))
			metrics.Count(ctx, "jobs.update.channels", id, len(channels.Channels))

			for i := range channels.Channels {
				err := backend.UpdateChannel(ctx, channels.Channels[i].ID, id, channels.Channels[i].Name, channels.Channels[i].Topic.Value, channels.Channels[i].Purpose.Value, channels.Channels[i].IsArchived, channels.Channels[i].IsPrivate, false)

				if err != nil {
					logger.Error(ctx, "jobs.update.channels", err.Error())
				} else {
					logger.Info(ctx, "jobs.update.channels", "channel=%s", channels.Channels[i].ID)
				}
			}

			nextCursor := channels.ResponseMetadata["next_cursor"]
			if nextCursor != "" {
				// there is more data, schedule its retrieval
				job.ScheduleJob(ctx, backend.BackgroundWorkQueue, fmt.Sprintf(types.JobsBaseURL+"/channels?id=%v&cursor=%v", id, nextCursor))
				logger.Info(ctx, "jobs.update.channels", "next=%s", nextCursor)
			}
		} else {
			// Slack API returned an error
			logger.Critical(ctx, "jobs.update.channels", "status=%s", channels.Error)
		}

	} else {
		logger.Error(ctx, "jobs.update.channels", err.Error())
	}
}

*/
