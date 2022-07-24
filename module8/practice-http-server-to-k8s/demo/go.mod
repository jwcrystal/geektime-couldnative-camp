module demo

go 1.18

require (
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	handlers v0.0.0
)

replace handlers => ./../pkg/handlers
