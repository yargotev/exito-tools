# Go package layout separates domains from surfaces

Exito Tools will use a Go package layout that separates application wiring, configuration, capability contracts, operational domains, interaction surfaces, and presenters. Domain packages must not import Cobra, Bubble Tea, or surface packages; surfaces adapt domain capabilities into CLI Commands and TUI Actions.
