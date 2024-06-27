/*
Generic Action Camera Adapter for 3DR Solo
Copyright (C) 2016 - 2017  Hugh Eaves

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/************************************************
* PART SELECTION
************************************************/

// Select the part you would like to display
part = "adapter";// [adapter:Adapter,retaining_clip:Retaining clip,back_weight:Back Weight,side_weight:Side Weight]

/********************
* MAIN PARAMETERS *
********************/

// No-name SJ4000 clone
// camera_length = 59.8;
// camera_width = 40.8;
// camera_height = 25.6;

// The values below are for the GitUp Git2P

/* [Adapter Parameters] */
// Length of camera body (side to side)
camera_length = 59.4;// [57:62]
// Width of camera body (top to bottom)
camera_width = 41.1;// [38:42]
// Height of camera body (front to back)
camera_height = 19.6;// [15:40]

// Adapter arm bump size (How far out the "bumps" protrude on the ends of the arms)
adapter_arm_bump_size = 2.5;// [1:5]

/* [Back Weight Parameters] */
// The length of the counterbalance weight on the back of the camera. Larger values make it heavier.
back_weight_length = 29;// [10:40]

// The offset of the weight on the back of the gimbal. Values greater than zero move the weight to the outer side of the gimbal
back_weight_offset = 0;// [0:20]

/* [Side Weight Parameters] */
// The thickness of the side weight, thicker is heavier
side_weight_thickness = 4;// [4:8]

// The length of the side weight, longer is heavier
side_weight_length = 10;// [9:17]

/********************
* OTHER PARAMETERS *
********************/
/* [Hidden] */

/* How much to shift the camera to the right to avoid the curved protrusion from the pitch motor housing */
camera_offset = 3.2;

/* Dimensions of the base of the camera adapter, so that it is a nice snug fit inside the 3DR gimbal mount */
base_height = 8;// (thickness) thick enough to clear the GoPro plug in the back
base_width = 42;
base_length = 60;

/* Dimensions of the small bracket that holds the camera adapter to the gimbal mount */
bracket_width = 9;
bracket_depth = 9.5;
bracket_thickness = 3;
bracket_radius = 2;

/* Typical values for M2 cap-head screws */
screw_thread_diameter = 2.5;
screw_head_diameter = 5;
screw_head_height = 3;

/* Dimensions of the "stud" that the 3DR balancing weights attach to */
stud_height = 2.5;
stud_diameter = 6.4;
stud_hole_diameter = 2;

bottom_lip_width = 4.5;
bottom_lip_height = 2.6;

// cutout for GoPro plug
plug_thickness = 5;
plug_width = 23;
plug_offset = 55.5 - (plug_thickness / 2);

// cutout for right side wall of "gimbal bracket"
side_wall_length = 31;
side_wall_height = 6.4;
side_wall_width = 3;
side_wall_offset = 59.4;

// height of hook on bottom left of "gimbal bracket"
hook_width = 17.6;
hook_height = 2.5;// distance from bottom
hook_length = 4;
hook_thickness = 1;

// width of protrusion on front of "pitch motor"
motor_protrusion_width = 20;

// thickness / length of camera retaining clips
clip_thickness = 1.6;
right_clip_opening_width = 26;
// left_clip_opening_width = motor_protrusion_width;
left_clip_opening_width = right_clip_opening_width;

x_midpoint = base_width / 2;
main_length = camera_length + camera_offset * 2;
camera_gap = (base_width - camera_width) / 2;

camera_height_adjustment = 0;

// more facets for smaller holes
$fn = 16;

module lip_cutout() {
  cube([bottom_lip_width, base_length, bottom_lip_height]);
  translate([0, 0, bottom_lip_height])
    mirror([1, -1, 0])
      wedge(base_length, bottom_lip_width, bottom_lip_width, false, false);
}


module side_wall_cutout() {
  translate([0, side_wall_offset, 0]) {
    cutout_base_height = side_wall_height - (side_wall_width / 2);
    cube([side_wall_length, side_wall_width, cutout_base_height + 0.01]);
    translate([0, 0, cutout_base_height])
      prism(side_wall_length, side_wall_width, side_wall_width / 2 - 0.01);
  }
}

module hook_cutout() {
  translate([x_midpoint - hook_width / 2, 0, hook_height]) {
    cube([hook_width, hook_length, hook_thickness]);
    translate([0, 0, hook_thickness])
      wedge(hook_width, hook_length, hook_length, false);
  }
}

module bracket_cutout() {
  bracket_clearance = 0.4;
  adj_bracket_width = bracket_width + bracket_clearance;
  adj_bracket_depth = bracket_depth + bracket_clearance;
  translate([base_width - bracket_depth, base_length + bracket_clearance / 2 - adj_bracket_width, 0]) {
    translate([-bracket_clearance, 0, base_height - bracket_thickness])
      cube([adj_bracket_depth, adj_bracket_width, bracket_thickness + 5]);
    translate([bracket_depth - bracket_thickness, 0, 0])
      cube([bracket_depth, adj_bracket_width, base_height]);
  }

}

module clip_corner(length) {

  cube([clip_thickness, length + clip_thickness, length]);
  wedge(clip_thickness, length + clip_thickness, length + clip_thickness, false, true);

}


module clip(length) {
  cube([base_width, clip_thickness, camera_height]);
  translate([-clip_thickness + camera_gap, 0, camera_height])
    cube([camera_width + clip_thickness * 2, clip_thickness, length * 2 + camera_height_adjustment]);
  translate([-clip_thickness + camera_gap, clip_thickness, camera_height + camera_height_adjustment])
    mirror([0, 1, -1])
      prism(camera_width + clip_thickness * 2, length * 2, length);
  //   cube([camera_width+clip_thickness*2, length + clip_thickness, clip_thickness ]);
  wedge(base_width, camera_offset - clip_thickness, camera_height + length * 2 + camera_height_adjustment, true, false);
  translate([0, 0, camera_height + adapter_arm_bump_size + camera_height_adjustment - length]) {
    translate([-clip_thickness + camera_gap, 0, 0])
      clip_corner(length);
    translate([base_width - camera_gap, 0, 0])
      clip_corner(length);
  }
}

module left_clip() {
  translate([0, camera_offset - clip_thickness, base_height]) {
    difference() {
      clip(adapter_arm_bump_size);
      translate([x_midpoint - left_clip_opening_width / 2, -clip_thickness, 2])
        cube([left_clip_opening_width, adapter_arm_bump_size + clip_thickness * 2, camera_height + adapter_arm_bump_size * 2]);
    }
  }
}

module right_clip() {
  translate([0, camera_length + camera_offset + clip_thickness, base_height]) {
    mirror([0, 1, 0])
      difference() {
        clip(adapter_arm_bump_size);
        translate([x_midpoint - right_clip_opening_width / 2, -clip_thickness, 2])
          cube([right_clip_opening_width, adapter_arm_bump_size + clip_thickness * 2, camera_height + adapter_arm_bump_size * 2]);
      }
  }
}

module prism(x, y, z, upside_down) {
  translate([0, y / 2, 0]) {
    wedge(x, y / 2, z, true, upside_down);
    wedge(x, y / 2, z, false, upside_down);
  }
}

module wedge(x, y, z, reverse, upside_down) {
  ly = reverse ? -y : y;
  lz = upside_down ? -z : z;
  polyhedron(points = [
    [0, 0, 0], 
    [x, 0, 0], 
    [x, ly, 0], 
    [0, ly, 0], 
    [0, 0, lz], 
    [x, 0, lz]
  ], faces = [
    [0, 1, 2, 3], 
    [5, 4, 3, 2], 
    [0, 4, 5, 1], 
    [0, 3, 4], 
    [5, 2, 1]
  ]);
}

module weight_stud() {
  difference() {
    cylinder(d = stud_diameter - 0.4, h = stud_height, $fn = 32);
    cylinder(d = stud_hole_diameter, h = stud_height, $fn = 16);
  }
}


/********************************************
* Top level function to render the adapter
********************************************/
module adapter() {

  difference() {
    cube([base_width, main_length, base_height]);

    lip_cutout();

    translate([base_width, 0, 0])
      mirror([-1, 0, 0])
        lip_cutout();

    side_wall_cutout();

    translate([x_midpoint - plug_width / 2, plug_offset, 0])
      cube([plug_width, plug_thickness, base_height]);

    hook_cutout();

    bracket_cutout();
  }

  left_clip();
  right_clip();

}

/********************************************
* Top level function to render the retaining clip
********************************************/
module retaining_clip() {
  bracket_inside_height = 11.4 + base_height - bracket_thickness;

  difference() {
    union() {
      hull() {
        cube([bracket_width, bracket_thickness, bracket_depth + bracket_thickness - bracket_radius]);
        cube([1, bracket_thickness, bracket_depth + bracket_thickness]);
        translate([bracket_width - bracket_radius, 0, bracket_depth + bracket_thickness - bracket_radius])
          mirror([0, -1, 1])
            cylinder(r = bracket_radius, h = bracket_thickness);
      }
      translate([0, bracket_inside_height + bracket_thickness, 0])
        cube([bracket_width, bracket_thickness, bracket_depth + bracket_thickness]);
      translate([0, bracket_thickness, 0])
        cube([bracket_width, bracket_inside_height, bracket_thickness * 2]);
    }

    translate([bracket_width - 3, bracket_thickness, bracket_thickness - 0.2])
      cube([2, bracket_inside_height, bracket_depth]);
    translate([bracket_width - 3.6, 0, bracket_depth + bracket_thickness - 2.8])
      mirror([0, -1, 1]) {
        translate([0, 0, -bracket_thickness / 2])
          cylinder(d = 4, h = bracket_thickness);
        translate([0, 0, bracket_thickness / 2])
          cylinder(d = 2, h = bracket_thickness);
      }
  }
}

/********************************************
* Top level function to render the back / counterbalance weight
********************************************/
module back_weight() {
  thickness = 2;
  length = 57.8;
  width = 34;
  weight_thickness = 100;

  intersection() {
    difference() {
      union() {
        cube([length, width, thickness]);
        translate([back_weight_offset, 0, -weight_thickness])
          cube([back_weight_length, width, weight_thickness]);
        translate([12, 0, thickness])
          cube([3, width, 3.4]);
        translate([33.6, 0, thickness])
          cube([3, width, 4]);
        translate([54.8, 8, thickness])
          cube([3, width - 8, 4]);

        translate([36, 0, 0]) {
          translate([4.5, 18, -2])
            cylinder(d = 5, h = 2);
          translate([16.5, 6.75, -2])
            cylinder(d = 5, h = 2);
          translate([13.5, 27, -2])
            cylinder(d = 5, h = 2);
        }

      }
      translate([36, 0, 0]) {
        translate([4.5, 18, -90])
          cylinder(d = screw_thread_diameter, h = 100);
        translate([16.5, 6.75, -90])
          cylinder(d = screw_thread_diameter, h = 100);
        translate([13.5, 27, -90])
          cylinder(d = screw_thread_diameter, h = 100);
      }

    }
    mirror([-1, 0, 1])
      translate([26, width / 2, 0]) {
        cylinder(d = 77, h = 14, $fn = 64);
        translate([0, 0, 14])
          cylinder(d = 70, h = 16, $fn = 64);
        translate([0, 0, 30])
          cylinder(d = 77, h = 100, $fn = 64);
      }
  }
}

module curve_cutout(beam_curve_diameter, beam_thickness) {
  difference() {
    cylinder(d = beam_curve_diameter, h = 100, $fn = 64);
    cylinder(d = beam_curve_diameter - (beam_thickness * 2), h = 100, $fn = 64);
  }
}

/********************************************
* Top level function to render the side weight
********************************************/
module side_weight() {

  side_weight_width = 2;

  beam_curve_diameter = 34;
  beam_thickness = 1.2;
  beam_height = 4;
  beam_width = 5.7;
  beam_length = beam_curve_diameter / 2;

  bracket_length = beam_length;
  bracket_width = beam_width + side_weight_width + beam_thickness;
  bracket_height = beam_height + side_weight_thickness;
  bracket_depth = 3;

  difference() {
    union() {
      difference() {
        cube([bracket_length, bracket_width, bracket_height]);

        translate([0, beam_curve_diameter / 2 + side_weight_width, side_weight_thickness])
          curve_cutout(beam_curve_diameter, beam_thickness);

        translate([bracket_length - screw_head_diameter / 2 - 1, bracket_width / 2 - 1, 0]) {

          cylinder(d = screw_head_diameter, h = screw_head_height);
          cylinder(d = screw_thread_diameter, h = bracket_height);
        }
      }
      translate([-side_weight_length, 0, 0])
        difference() {
          cube([side_weight_length, bracket_width, bracket_height]);

          translate([0, side_weight_width, side_weight_thickness])
            cube([side_weight_length, beam_thickness, 100]);
        }
    }
  }
}

// Call the top level rendering function based on "part"
if (part == "adapter") {
  adapter();
} else if (part == "back_weight") {
  back_weight();
} else if (part == "retaining_clip") {
  retaining_clip();
} else if (part == "side_weight") {
  side_weight();
  translate([0, 20, 0])
    mirror([1, 0, 0])
      side_weight();
}
