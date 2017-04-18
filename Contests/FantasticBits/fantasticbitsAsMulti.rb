STDOUT.sync = true
@my_team_id = gets.to_i # if 0 you need to score on the right of the map, if 1 you need to score on the left

@my_team_id == 0 ? $GOAL = [16000, 3750] : $GOAL = [0, 3750]

class Game
  def initialize(mScore, mMagic, oScore, oMagic)
    @myS = mScore
    @myM = mMagic
    @theirS = oScore
    @theirM = oMagic
    @myW = []
    @theirW = []
    @snaff = []
    @blud = []
  end
  class << self
    attr_accessor :myS, :myM, :theirS, :theirM, :myW, :theirW, :snaff, :blud
  end
end

loop do
    my_score, my_magic = gets.split(" ").collect {|x| x.to_i}
    opponent_score, opponent_magic = gets.split(" ").collect {|x| x.to_i}
    game = Game.new(my_score, my_magic, opponent_score, opponent_magic)
    entities = gets.to_i # number of entities still in game
    entities.times do
        # entity_id: entity identifier
        # type: "WIZARD", "OPPONENT_WIZARD" or "SNAFFLE" (or "BLUDGER" after first league)
        # x: position
        # y: position
        # vx: velocity
        # vy: velocity
        # state: 1 if the wizard is holding a Snaffle, 0 otherwise
        entity_id, type, x, y, vx, vy, state = gets.split(" ")
        entity = [entity_id.to_i, x.to_i, y.to_i, vx.to_i, vy.to_i, state.to_i]
        (game.myW.push(entity) if type == "WIZARD") ||
        (game.theirW.push(entity) if type == "OPPONENT_WIZARD") ||
        (game.snaff.push(entity)) if type == "SNAFFLE" ||
        (game.blud.push(entity)) if type == "BLUDGER"
    end

    game.myW do |wiz|

        if wiz[5] == 1 then
          printf("THROW", $GOAL[0], $GOAL[1], 500)
        else
          closest_coordinates = [wiz[1], wiz[2]]
          min_distance = 100000
          game.snaff do |snaff|
            distance = Math.sqrt((wiz[1] - snaff[1])**2 + (wiz[2] - snaff[2])**2)
            if min_distance > distance then
              min_distance = distance
              closest_coordinates = [snaff[1], snaff[2]]
            end
          end
          printf("MOVE", closest_coordinates[1], closest_coordinates[2], "150")
        end
        # To debug: STDERR.puts "Debug messages..."
    end
end
