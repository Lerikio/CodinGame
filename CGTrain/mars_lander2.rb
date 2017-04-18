STDOUT.sync = true # DO NOT REMOVE

Point = Struct.new(:x, :y)

def compute_intersection(a, b, c, d)
  p = a
  r = Point.new(b.x-p.x, b.y-p.y)
  q = c
  s = s = Point.new(d.x - q.x, d.y - q.y)
  if r.x*s.y-r.y*s.x != 0
    t = ((q.x-p.x)*s.y - (q.y-p.y)*s.x)/(r.x*s.y-r.y*s.x)
    u = ((q.x-p.x)*r.y - (q.y-p.y)*r.x)/(r.x*s.y-r.y*s.x)
    if 0 < u && u < 1 && 0 < t && t< 1
      return Point.new(p.x + t*r.x, p.y + t*r.y)
    end
  end
  return false
end

def compute_intersections(start, stop, surface, forward)
  intersections = []
  for i in 0..(surface.length-2)
    if forward then
      current_surface = surface[i]
      next_surface = surface[i+1]
    else
      current_surface = surface[-i]
      next_surface = surface[-(i+1)]
    end
    intersection = compute_intersection(current_surface, next_surface, start, stop)
    if intersection
      distance = Math.sqrt((start.x - intersection.x)**2 + (start.y - intersection.y)**2)
      intersections << [i, distance]
    end
  end
  return intersections
end

def compute_trajectory(trajectory, landing, surface)

  current_start = trajectory[-1]

  # Détermine dans quel sens parcourir la surface
  current_start.x < landing.x ? forward = true : forward = false

  # Détermine les intersections entre la surface et la trajectoire directe entre le point actuel et l'aterrissage
  intersections = compute_intersections(current_start, landing, surface, forward)
  STDERR.puts intersections
  if intersections.length == 0
    # La trajectoire est bonne !
    return trajectory << landing
  else
    # Détermine l'intersection la plus proche
    closest = -1
    shortest_distance = 1000000
    intersections.each do |index, dist|
      if dist < shortest_distance
        closest = index
        shortest_distance = dist
      end
    end

    # Trouve le chemin libre le plus proche
    next_point = Point.new(0, 0)
    while closest < surface.length-3 do
      if forward then
        current_surface = surface[closest]
        next_surface = surface[closest+1]
      else
        current_surface = surface[-closest]
        next_surface = surface[-(closest+1)]
      end

      length_surface = Math.sqrt((current_surface.x - next_surface.x)**2 + (current_surface.y - next_surface.y)**2)
      test_point = Point.new(next_surface.x + 3 * (next_surface.x - current_surface.x) / length_surface,
                             next_surface.y + 3 * (next_surface.y - current_surface.y) / length_surface)

      if !compute_intersection(next_surface, surface[closest+2], current_start, test_point)
        next_point = test_point
        break
      else
        # Incrémente le compteur
        closest += 1
      end
    end
    STDERR.puts next_point.x.to_s + " " + next_point.y.to_s
    if next_point.x != 0 || next_point.y != 0
      trajectory << next_point
      return compute_trajectory(trajectory, landing, surface)
    else
      STDERR.puts "No trajectory found..."
      return false
    end
  end
end

@surface = []
@landing = [0, 0, 0]

@surface_n = gets.to_i # the number of points used to draw the surface of Mars.
@surface_n.times do
    # land_x: X coordinate of a surface point. (0 to 6999)
    # land_y: Y coordinate of a surface point. By linking all the points together in a sequential fashion, you form the surface of Mars.
    land_x, land_y = gets.split(" ").collect {|x| x.to_i}
    STDERR.puts land_x.to_s + " " + land_y.to_s
    @landing = [@surface[-1][0], land_x, land_y] if @surface.length > 0 && land_y == @surface[-1][1]
    @surface << Point.new(land_x, land_y)
end

# game loop
loop do
    # h_speed: the horizontal speed (in m/s), can be negative.
    # v_speed: the vertical speed (in m/s), can be negative.
    # fuel: the quantity of remaining fuel in liters.
    # rotate: the rotation angle in degrees (-90 to 90).
    # power: the thrust power (0 to 4).
    x, y, h_speed, v_speed, fuel, rotate, power = gets.split(" ").collect {|x| x.to_i}

    # Search for an intersection between the direct line toward the landing and the ground
    start = Point.new(x, y)
    landing = Point.new(@landing[0] + (@landing[1] - @landing[0])/2, @landing[2] + 1)
    trajectory = compute_trajectory([start], landing, @surface)
    STDERR.puts trajectory

    closest_solution = [0, 0]
    closest_distance = 10000000

    for thrust in 0..4 do
      for angle in -90..90 do
        new_position = Point.new(start.x + h_speed - Math.sin(angle)*thrust, start.y + v_speed + Math.cos(angle)*thrust + 3.711)
        intersections = compute_intersections(start, new_position, @surface, true)
        if intersections.length == 0
          trajectory_length = Math.sqrt((start.x - trajectory[1].x)**2 + (start.y - trajectory[1].y)**2)
          ab = Point.new(trajectory[1].x - start.x, trajectory[1].y - start.y)
          ac = Point.new(new_position.x - start.x, new_position.y - start.y)
          dist_to_trajectory = (1 - (ab.x * ac.x + ab.y * ac.y)/trajectory_length**2).abs
          # dist_to_destination = Math.sqrt( (new_position.x - trajectory[1].x)**2 +  (new_position.y - trajectory[1].y)**2)
          if dist_to_trajectory < closest_distance
            STDERR.puts dist_to_trajectory.to_s + " < " + closest_distance.to_s
            closest_distance = dist_to_trajectory
            closest_solution = [thrust, angle]
          end
        end
      end
    end

    # To debug: STDERR.puts "Debug messages..."
    puts closest_solution[1].to_s + " " + closest_solution[0].to_s
end
