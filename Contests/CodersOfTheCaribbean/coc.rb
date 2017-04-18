STDOUT.sync = true # DO NOT REMOVE

Point = Struct.new(:x, :y)

Barrel = Struct.new(:x, :y)
Ship = Struct.new(:x, :y, :d, :s, :owner, :attacked)

def next_position(ship)
  result = Point.new(ship.x, ship.y)
  if ship.d == 0 then
    delta = Point.new(1, 0)
  elsif ship.d == 1 then
    delta = Point.new(1, -1)
  elsif ship.d == 2 then
    delta = Point.new(0, -1)
  elsif ship.d == 3 then
    delta = Point.new(-1, 0)
  elsif ship.d == 4 then
    delta = Point.new(0, 1)
  elsif ship.d == 5 then
    delta = Point.new(1, 1)
  end

  ship.s.times do
    result.x = result.x + delta.x if result.x + delta.x < 23 && result.x + delta.x >= 0
    result.y = result.y + delta.y if result.y + delta.y < 21 && result.y + delta.y >= 0
  end

  return result
end

# game loop
loop do
    me = []
    not_me = []
    barrels = []
    my_ship_count = gets.to_i # the number of remaining ships
    entity_count = gets.to_i # the number of entities (e.g. ships, mines or cannonballs)
    entity_count.times do
        entity_id, entity_type, x, y, arg_1, arg_2, arg_3, arg_4 = gets.split(" ")
        entity_id = entity_id.to_i
        x = x.to_i
        y = y.to_i
        arg_1 = arg_1.to_i
        arg_2 = arg_2.to_i
        arg_3 = arg_3.to_i
        arg_4 = arg_4.to_i
        if entity_type == "SHIP" && arg_4 == 1 then
            me << Ship.new(x, y, arg_1, arg_2, arg_4, false)
        elsif entity_type == "SHIP" && arg_4 == 0 then
            not_me << Ship.new(x, y, arg_1, arg_2, arg_4, false)
        elsif entity_type == "BARREL" then
            barrels << Barrel.new(x, y)
        end
    end
    me.each do |ship|
      target = Point.new(-1, -1)
      not_me.each do |them|
        distance_ship = ((ship.x - them.x).abs + (ship.x + ship.y - them.x - them.y).abs + (ship.y - them.y).abs) / 2
        if distance_ship < 10 && !ship.attacked then
          ship.attacked  = true
          target = next_position(them)
        end
      end
        if target.x >= 0 then
          puts "FIRE " + target.x.to_s + " " + target.y.to_s
        else
          ship.attacked  = false
          closest_barrel = Barrel.new(0, 0)
          min = 1000000
          barrels.each do |b|
              distance = Math.sqrt((ship.x - b.x)**2 + (ship.y - b.y)**2)
              if distance < min
                  closest_barrel = b
                  min = distance
              end
          end

          # To debug: STDERR.puts "Debug messages..."

          printf("MOVE " + closest_barrel.x.to_s + " " + closest_barrel.y.to_s + "\n") # Any valid action, such as "WAIT" or "MOVE x y"
        end
    end
end
