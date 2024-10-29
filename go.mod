module github.com/kinta-mti/mobbe

go 1.20

require (
	github.com/kinta-mti/mobbe/config v0.0.0-20241029043912-808e8b5b0011
	github.com/kinta-mti/mobbe/db v0.0.0-20241029081110-42cced18400f
	github.com/kinta-mti/mobbe/ypg v0.0.0-20241029064354-6151091fe6c5
)

replace (
	github.com/kinta-mti/mobbe/config => ./config
	github.com/kinta-mti/mobbe/ypg => ./ypg
	github.com/kinta-mti/mobbe/db => ./db
)
