package applications

type (
	AppsTable map[string][]AppDetails

	AppDetails struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Icon        string `json:"icon"`
	}
)

var ExampledtApps = AppsTable{
	"SelfHosting": []AppDetails{
		{
			Name:        "Nextcloud",
			Description: "Nextcloud is a self-hosted cloud storage solution that allows you to store, sync, and share files securely.",
			Link:        "https://nextcloud.com",
			Icon:        "nextcloud",
		},
		{
			Name:        "Plex",
			Description: "Plex is a media server software that lets you organize and stream your media content to multiple devices.",
			Link:        "https://www.plex.tv",
			Icon:        "plex",
		},
		{
			Name:        "Home Assistant",
			Description: "Home Assistant is an open-source home automation platform that focuses on local control and privacy for smart home devices.",
			Link:        "https://www.home-assistant.io",
			Icon:        "home-assistant",
		},
	},
	"OtherTools": []AppDetails{
		{
			Name:        "Grafana",
			Description: "Grafana is an open-source analytics and monitoring platform that allows you to query, visualize, and alert on the metrics in real-time.",
			Link:        "https://grafana.com",
			Icon:        "grafana",
		},
		{
			Name:        "Jellyfin",
			Description: "Jellyfin is a free and open-source media server software that allows you to stream videos, music, and other media content to various devices.",
			Link:        "https://jellyfin.org",
			Icon:        "jellyfin",
		},
	},
}
