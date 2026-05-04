package models

type MBuild struct {
	Dockerfile string            `json:"dockerfile,omitempty"`
	Context    string            `json:"context,omitempty"`
	Args       map[string]string `json:"args,omitempty"`
}

type MPortAttributes struct {
	Label           string `json:"label,omitempty"`
	OnAutoForward   string `json:"onAutoForward,omitempty"` // "notify", "openBrowser", "openPreview", "silent", "ignore"
	ElevateIfNeeded bool   `json:"elevateIfNeeded,omitempty"`
}

type MCustomizations struct {
	VScode    MVScode    `json:"vscode"`
	JetBrains MJetBrains `json:"jetbrains"`
}

type MDevContainer struct {
	// Nome de exibição do container
	Name string `json:"name,omitempty"`

	// --- CONFIGURAÇÃO DA IMAGEM/DOCKER ---
	Image string `json:"image,omitempty"`
	Build MBuild `json:"build"`

	// --- REDE E PORTAS ---
	// Encaminha portas do container para o host local
	ForwardPorts []int `json:"forwardPorts,omitempty"`
	// Atributos específicos para as portas
	PortsAttributes map[string]MPortAttributes `json:"portsAttributes,omitempty"`

	// --- CICLO DE VIDA E EXECUÇÃO ---
	// Comandos executados em diferentes estágios da criação
	OnCreateCommand      string `json:"onCreateCommand,omitempty"`
	UpdateContentCommand string `json:"updateContentCommand,omitempty"`
	PostCreateCommand    string `json:"postCreateCommand,omitempty"`
	PostStartCommand     string `json:"postStartCommand,omitempty"`

	// --- PERMISSÕES E ACESSO ---
	RemoteUser    string `json:"remoteUser,omitempty"`
	ContainerUser string `json:"containerUser,omitempty"`
	Privileged    bool   `json:"privileged,omitempty"`

	// --- VARIÁVEIS DE AMBIENTE ---
	ContainerEnv map[string]string `json:"containerEnv,omitempty"`

	// Others
	WorkspaceFolder string          `json:"workspaceFolder,omitempty"`
	WorkspaceMount  string          `json:"workspaceMount,omitempty"`
	Customizations  MCustomizations `json:"customizations"`
}
