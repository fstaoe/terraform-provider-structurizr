package model

// Workspace represents a workspace configured in the structurizr
type Workspace struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	APIKey       string `json:"apiKey"`
	APISecret    string `json:"apiSecret"`
	PublicURL    string `json:"publicUrl"`
	PrivateURL   string `json:"privateUrl"`
	ShareableURL string `json:"shareableUrl"`
}

// Workspaces is the response body for any CRU methods
type Workspaces struct {
	Workspaces []*Workspace `json:"workspaces"`
}

// FindByID returns a workspace by its ID
func (w *Workspaces) FindByID(id any) *Workspace {
	for _, workspace := range w.Workspaces {
		if workspace.ID == id {
			return workspace
		}
	}
	return nil
}
