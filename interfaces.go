package main

type (
	// Parses photos of some hotel
	PhotoParser interface {
		ParseHotelsByCountry(city string)
	}

	// Generates frames for future movie
	FrameGenerator interface {
		GenerateFrames(photosPath string) error //returns path of the generated frames or error
	}

	//Creats movie from generated frames
	MovieMaker interface {
		MakeMovie(movieName string) (string, error) //returns path of the created movie or error
	}

	//Uploads movie on youtube
	YoutubeUploader interface {
		UploadMovie(moviePath string) error // returns error
	}
)
