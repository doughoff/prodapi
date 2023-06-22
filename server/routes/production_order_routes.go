package routes

// production order should be like this :
/*
	# Production Order
	  - id
 	  - recipe //etc, other fields.

	  -> Cycles
		- All production Order Cycles (based on recipe/targetAmount)

	  -> Movements
		- All stock movements "from the cycles" (use the stockMovement Query to bind things together)


	// TODO: Probably is needed to refactor a little the stock movements part to make something "reusable" for this.

*/
