var Point = function  (x, y) {
	this.x = x;
	this.y = y;
};

var Goal = function (center, mine) {
	this.x = center;
	this.mine = mine;
	this.top = 3750 - 2000 + 300;
	this.bottom = 3750 + 2000 - 300;
};

var Moving = function (x, y, vx, vy, radius) {
	this.center = new Point(x, y);
	this.speed = new Point(vx, vy);
	this.radius = radius;
};

Moving.prototype.distanceTo = function(other) {
	var dist_x, dist_y, distance;
	dist_x = abs(this.center.x - other.center.x);
	dist_y = abs(this.center.y - other.center.y);
	distance = sqrt(dist_x * dist_x + dist_y * dist_y);
	return distance;
}

var Snaffle = function (x, y, vx, vy) {
	Moving.call(this, x, y, vx, vy, 150);
}
Snaffle.prototype = Object.create(Moving.Prototype);
Snaffle.prototype.constructor = Snaffle;

var Wizard = function (x, y, vx, vy, owner, snaffle) {
	Moving.call(this, x, y, vx, vy, 400);
	this.mine = owner;
	this.snaffleing = snaffle;
}
Wizard.prototype = Object.create(Moving.Prototype);
Wizard.prototype.constructor = Wizard;

Wizard.prototype.findClosestSnaffle = function (snaffles, ignore) {
	var closestIndex = 0;
	var distance = 30000000;
	for (var i = snaffles.length - 1; i >= 0; i--) {
		if (i != ignore) {
			var curDistance = this.distanceTo(snaffles[i]);
			if (curDistance < distance) {
				closestIndex = i;
				distance = curDistance;
			}
		}
	}
}

/**
 * Grab Snaffles and try to throw them through the opponent's goal!
 * Move towards a Snaffle and use your team id to determine where you need to throw it.
 **/

var myTeamId = parseInt(readline()); // if 0 you need to score on the right of the map, if 1 you need to score on the left

var myGoal = new Goal(16000, true);
var theirGoal = new Goal(0, false);

if (myTeamId == 0)
{
    myGoal = new Goal(0, true);
    theirGoal = new Goal(16000, false);
}

// game loop
while (true) {
	var myWizards = [];
	var theirWizards = [];
	var snaffles = [];

    var entities = parseInt(readline()); // number of entities still in game
    
    for (var i = 0; i < entities; i++) {
        var inputs = readline().split(' ');
        var entityId = parseInt(inputs[0]); // entity identifier
        var entityType = inputs[1]; // "WIZARD", "OPPONENT_WIZARD" or "SNAFFLE" (or "BLUDGER" after first league)
        var x = parseInt(inputs[2]); // position
        var y = parseInt(inputs[3]); // position
        var vx = parseInt(inputs[4]); // velocity
        var vy = parseInt(inputs[5]); // velocity
        var state = parseInt(inputs[6]); // 1 if the wizard is holding a Snaffle, 0 otherwise
    
    	if (entityType == "WIZARD"){
    		myWizards.push(new Wizard(x, y, vx, vy, true, state))
    	} else if (entityType == "OPPONENT_WIZARD") {
    		theirWizards.push(new Wizard(x, y, vx, vy, true, state))
    	} else if (entityType == "SNAFFLE") {
    		snaffles.push(new Snaffle(x, y, vx, vy))
    	}

    }
    
    var first_closest = -1;

    for (var i = 0; i < 2; i++) {


    	if (myWizards[i].snaffleing) {
    		print('THROW ' + theirGoal.x + ' 3750 500');
    	} else {
    		first_closest = myWizards[i].findClosestSnaffle(snaffles, first_closest);
    		print('MOVE ' + snaffles[first_closest].center.x + ' ' + snaffles[first_closest].center.y + ' 150' );
    	}
        // Write an action using print()
        // To debug: printErr('Debug messages...');
    }
}