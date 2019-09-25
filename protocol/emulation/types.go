package emulation


// Screen orientation.
type ScreenOrientation struct {
	
	// Orientation type.
	// Possible value: portraitPrimary,portraitSecondary,landscapePrimary,landscapeSecondary,
	Type	string	`json:"type"`
	
	// Orientation angle.
	
	Angle	int	`json:"angle"`
	
}	

// advance: If the scheduler runs out of immediate work, the virtual time base may fast forward to
	// allow the next delayed task (if any) to run; pause: The virtual time base may not advance;
	// pauseIfNetworkFetchesPending: The virtual time base may not advance if there are any pending
	// resource fetches.
type VirtualTimePolicy string	

