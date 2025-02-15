include <../common/shapes.scad>
include <../common/util.scad>

$fs = 0.6;
$fa = 6;

opt_active_cooling = false;
opt_rpi_5 = true;
opt_mounting_tabs = false;

case_thickness = 1.6;
width_wall_gap = 1;
length_wall_gap = 1.0;
case_top_gap = 1;

flange_height = 1.2;
flange_width = 1.2;
flange_gap = 0.2;
post_gap = flange_gap + 0.2;

grid_hole_width = 3;
grid_bar_width = 1.0;

tab_width = 10;
tab_thickness = 3;
tab_screw_diameter = 4;

screw_head_diameter = 5.4;
screw_head_thickness = 2;
screw_thread_major_diameter = 3;// Note: pcb has 2.7mm hole on mechanical drawing
screw_thread_minor_diameter = 2.4;

pcb_length = 85;
pcb_width = 56;
pcb_corner_radius = 3.0;
pcb_screw_offset = 3.5;
pcb_screw_length_offset = 58 + pcb_screw_offset;
pcb_thickness = 1.6;
pcb_back_clearance = 3.6;

pcb_standoff_diameter = 6.2;// Note: 6mm wide pad on mechanical drawing

jack_length = 20;
jack_overhang = 4;
jack_width_gap = 0.4;
jack_height_gap = 0.2;

sdcard_opening_width = 12;
sdcard_opening_height = pcb_back_clearance + case_thickness;
sdcard_opening_length = 4;
sdcard_opening_pos = [-jack_overhang, pcb_width / 2, 0];
sdcard_thickness = 1.4;

usb_c_jack_width = 9;
usb_c_jack_height = 3.2;
usb_c_jack_offset=pcb_screw_offset + 7.7;
usb_c_jack_pos = [usb_c_jack_offset, -jack_overhang, 0];

hdmi_jack_width = 7.2;
hdmi_jack_height = 3.5;
hdmi_jack_1_offset=usb_c_jack_offset + 14.8;
hdmi_jack_1_pos = [hdmi_jack_1_offset, -jack_overhang, 0];
hdmi_jack_2_offset=hdmi_jack_1_offset + 13.5;
hdmi_jack_2_pos = [hdmi_jack_2_offset, -jack_overhang, 0];

audio_jack_diameter = 6;
audio_jack_height = audio_jack_diameter;
audio_jack_width = 7;
audio_jack_vertical_offset = 0.2;
audio_jack_offset=hdmi_jack_2_offset+14.5;
audio_jack_pos = [audio_jack_offset, -jack_overhang, 0];
audio_jack_gap = 0.6;

usb_a_jack_width = 13.2;
usb_a_jack_height = 14.8;
usb_a_jack_vertical_offset = 1.4;
usb_a_jack_1_pos = [pcb_length - jack_length + jack_overhang, opt_rpi_5 ? 29.1 : 9, 0];
usb_a_jack_2_pos = [pcb_length - jack_length + jack_overhang, opt_rpi_5 ? 47: 27, 0];

ethernet_jack_pos = [pcb_length - jack_length + jack_overhang, opt_rpi_5 ? 10.2 : 45.75, 0];
ethernet_jack_width = 16.1;
ethernet_jack_height = 13.5;
ethernet_jack_vertical_offset = 0.6;

// Raspberry Pi 5 Power Button
power_button_pos = [-jack_overhang, 18.4, 0];
power_button_vertical_offset = 0.4;
power_button_height = 3;
power_button_width = 4;
power_button_length = 10;

fan_width = 30;
fan_diameter = 28.6;
fan_screw_spacing = 24;
fan_screw_hole_diameter = 3.2;
fan_retainer_height = 4;
fan_retainer_gap = 0.2;

jack_midpoint = max(audio_jack_diameter / 2 + audio_jack_vertical_offset, hdmi_jack_height / 2, usb_c_jack_height / 2);

jack_max_height = max(ethernet_jack_height + ethernet_jack_vertical_offset + +jack_height_gap / 2, usb_a_jack_height + usb_a_jack_vertical_offset + jack_height_gap / 2);

bottom_height = pcb_back_clearance + case_thickness + pcb_thickness + jack_midpoint;

top_height = case_thickness + jack_max_height - jack_midpoint + case_top_gap;

outer_length = pcb_length + (case_thickness + length_wall_gap) * 2;
outer_width = pcb_width + (case_thickness + width_wall_gap) * 2;
outer_radius = pcb_corner_radius + case_thickness + width_wall_gap;
inner_length = outer_length - case_thickness * 2;
inner_width = outer_width - case_thickness * 2;
inner_radius = pcb_corner_radius + width_wall_gap;

pcb_screw_positions = [
  [pcb_screw_offset, pcb_screw_offset, 0],
  [pcb_screw_offset, pcb_width - pcb_screw_offset, 0],
  [pcb_screw_length_offset, pcb_screw_offset, 0],
  [pcb_screw_length_offset, pcb_width - pcb_screw_offset, 0]
];

module hdmi_jack() {
  translate([0, 0, -jack_height_gap / 2])
    centered_cube([hdmi_jack_width + jack_width_gap, jack_length, hdmi_jack_height + jack_height_gap], [1, 0, 0]);
}

module usb_c_jack() {
  translate([0, 0, -jack_height_gap / 2])
    centered_cube([usb_c_jack_width + jack_width_gap, jack_length, hdmi_jack_height + jack_height_gap], [1, 0, 0]);
}

module audio_jack() {
  translate([0, jack_overhang - audio_jack_gap, 0])
    centered_cube([audio_jack_width + audio_jack_gap, jack_length - jack_overhang + audio_jack_gap, audio_jack_height + audio_jack_gap], [1, 0, 0]);

  translate([0, 0, -audio_jack_gap / 2 + audio_jack_vertical_offset])
    translate([0, jack_length, audio_jack_diameter / 2])
      rotate([90, 0, 0])
        cylinder(d = audio_jack_diameter + audio_jack_gap, h = jack_length);
}

module usb_a_jack() {
  translate([0, 0, -jack_height_gap / 2 + usb_a_jack_vertical_offset])
    centered_cube([jack_length, usb_a_jack_width + jack_width_gap, usb_a_jack_height + jack_height_gap], [0, 1, 0]);
}

module ethernet_jack() {
  translate([0, 0, -jack_height_gap / 2 + ethernet_jack_vertical_offset])
    centered_cube([jack_length, ethernet_jack_width + jack_width_gap, ethernet_jack_height + jack_height_gap], [0, 1, 0]);
}

module sdcard_opening() {
  translate([0, 0, -pcb_thickness - sdcard_opening_height * 2])
    centered_rounded_cube([sdcard_opening_length + jack_width_gap, sdcard_opening_width, sdcard_opening_height * 2], [0, 1, 0], r = 1);
}

module power_button() {
  translate([0, 0, power_button_vertical_offset]) {
    centered_cube([power_button_length, power_button_width,power_button_height],[0,1,0]);
    translate([power_button_length - 7.6, 0, 0])
   centered_cube([2.2, power_button_width + 3.2,power_button_height],[0,1,0]);
  }
}

module connectors() {
  translate(usb_c_jack_pos)
    usb_c_jack();
  translate(hdmi_jack_1_pos)
    hdmi_jack();
  translate(hdmi_jack_2_pos)
    hdmi_jack();

  if (!opt_rpi_5)
    translate(audio_jack_pos)
      audio_jack();

  translate(usb_a_jack_1_pos)
    usb_a_jack();
  translate(usb_a_jack_2_pos)
    usb_a_jack();
  translate(ethernet_jack_pos)
    ethernet_jack();
  if (opt_rpi_5)
    translate(power_button_pos)
      power_button();
}

module board() {
  difference() {
    union() {
      centered_rounded_cube([pcb_length, pcb_width, pcb_thickness], axes = [0, 0, 0], r = pcb_corner_radius, top = false, bottom = false);
    }
    for(pcb_screw_position = pcb_screw_positions)
      translate(pcb_screw_position)
        cylinder(h = case_thickness + pcb_back_clearance, d = screw_thread_major_diameter);
  }
  translate([0, 0, pcb_thickness])
    connectors();
}

module top_case_body() {

  difference() {
    union() {
      rounded_cube_2([outer_length, outer_width, top_height], side_r = outer_radius, bottom_r = 0.001, top_r = 1);

      translate([case_thickness + flange_gap, case_thickness + flange_gap, -flange_height])
        rounded_cube([inner_length - flange_gap * 2, inner_width - flange_gap * 2, flange_height + 0.01], r = inner_radius, top = false, bottom = false);
    }

    hull() {
      translate([case_thickness, case_thickness, flange_height + flange_width])
        // Note: add an extra 4mm of thickness of the USB/Ethernet end of the case
        rounded_cube([inner_length - 4, inner_width, top_height - case_thickness - flange_height * 2], r = inner_radius, top = false, bottom = false);

      translate([case_thickness + flange_gap + flange_width, case_thickness + flange_gap + flange_width, flange_height])
        rounded_cube([inner_length - flange_gap * 2 - flange_width * 2, inner_width - flange_gap * 2 - flange_width * 2, flange_height * 2 + 0.01], r = inner_radius - flange_width - flange_gap, top = false, bottom = false);
    }


    translate([case_thickness + flange_gap + flange_width, case_thickness + flange_gap + flange_width, -flange_height - 0.01])
      rounded_cube([inner_length - flange_gap * 2 - flange_width * 2, inner_width - flange_gap * 2 - flange_width * 2, flange_height * 2 + 0.02], r = inner_radius - flange_width - flange_gap, top = false, bottom = false);

  }

}

module fan_retainer() {
  difference() {
    centered_rounded_cube([fan_width + case_thickness * 2 + fan_retainer_gap * 2, fan_width + case_thickness * 2 + fan_retainer_gap * 2, fan_retainer_height], top = false, bottom = false);
    centered_rounded_cube([fan_width + fan_retainer_gap * 2, fan_width + fan_retainer_gap * 2, fan_retainer_height], top = false, bottom = false);
  }
}

module top() {
  difference() {
    union() {
      top_case_body();

      // screw posts
      translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, -jack_midpoint + post_gap]) {
        for(pcb_screw_position = pcb_screw_positions)
          translate(pcb_screw_position)
            cylinder(h = jack_midpoint + top_height - case_thickness, d = pcb_standoff_diameter);
      }
    }

    // screw post holes
    translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, -jack_midpoint + post_gap]) {
      for(pcb_screw_position = pcb_screw_positions)
        translate(pcb_screw_position)
          cylinder(h = 8, d = screw_thread_minor_diameter);
    }

    // connectors
    translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, -jack_midpoint]) {
      connectors();
    }

    if(opt_active_cooling) {
      translate([29, 33, 0]) {
        for(x = [-fan_screw_spacing / 2, fan_screw_spacing / 2]) {
          for(y = [-fan_screw_spacing / 2, fan_screw_spacing / 2]) {
            translate([x, y, 0])
              cylinder(d = fan_screw_hole_diameter, h = 50);
          }
        }

        cylinder(d = fan_diameter, h = 50);
      }
    }

    if(opt_active_cooling) {
      translate([63, outer_width / 2, 0])
        vent_grid(7, 10);
      translate([7, outer_width / 2, 0])
        vent_grid(2, 10);
    }

    if(!opt_active_cooling) {
      translate([outer_length / 2 - 38, outer_width / 2, 0])
        vent_grid(1, 11);
      translate([outer_length / 2 - 10, outer_width / 2, 0])
        vent_grid(13, 13);
      translate([outer_length / 2 + 22, outer_width / 2, 0])
        vent_grid(3, 11);
    }
  }

  if(opt_active_cooling) translate([29, 33, top_height - case_thickness - fan_retainer_height])
    fan_retainer();
}

module bottom_case_body() {
  difference() {
    union() {
      rounded_cube_2([outer_length, outer_width, bottom_height], side_r = outer_radius, top_r = 0.001, bottom_r = 1);
      if (opt_mounting_tabs)
        mounting_tabs();
    }

    translate([case_thickness, case_thickness, case_thickness])
      rounded_cube([inner_length, inner_width, bottom_height], r = inner_radius, top = false, bottom = false);

    if (opt_mounting_tabs)
      mounting_tab_holes();
  }
}

module bottom() {
  difference() {
    union() {
      bottom_case_body();

      // standoffs
      translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, case_thickness - 0.01]) {
        for(pcb_screw_position = pcb_screw_positions)
          translate(pcb_screw_position)
            cylinder(h = pcb_back_clearance, d1 = screw_head_diameter + 3.2, d2 = pcb_standoff_diameter);
      }

      translate(sdcard_opening_pos)
        translate([jack_overhang, case_thickness + width_wall_gap, case_thickness])
          centered_cube([sdcard_opening_length + 2, sdcard_opening_width + 2, pcb_back_clearance - sdcard_thickness], [0, 1, 0]);

    }

    translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, 0]) {

    // screw holes
      for(pcb_screw_position = pcb_screw_positions)
        translate(pcb_screw_position) {
          cylinder(h = pcb_back_clearance + case_thickness, d = screw_thread_major_diameter);
          cylinder(h = screw_head_thickness, d = screw_head_diameter);
        }

        // connector holes
      translate([0, 0, pcb_back_clearance + case_thickness + pcb_thickness])
        connectors();

      translate([0, 0, pcb_back_clearance + case_thickness + pcb_thickness])
        translate(sdcard_opening_pos)
          sdcard_opening();
    }

    // ventilation holes
    translate([45, outer_width / 2, 0]) {
      vent_grid(18, 9);
      translate([-10, 0, 0]) {
        vent_grid(11, 13);
      }

    }
  }
}

module mounting_tabs() {
  translate([tab_width / 2, outer_width / 2, 0]) {
    hull() {
      translate([0, 36, 0])
        rounded_cylinder(d = tab_width, h = tab_thickness, r = 1, top = true, bottom = true);
      translate([0, -36, 0])
        rounded_cylinder(d = tab_width, h = tab_thickness, r = 1, top = true, bottom = true);
    }
  }
  translate([outer_length - tab_width / 2, outer_width / 2, 0]) {
    hull() {
      translate([0, 36, 0])
        rounded_cylinder(d = tab_width, h = tab_thickness, r = 1, top = true, bottom = true);
      translate([0, -36, 0])
        rounded_cylinder(d = tab_width, h = tab_thickness, r = 1, top = true, bottom = true);
    }
  }
}

module mounting_tab_holes() {
  translate([tab_width / 2, outer_width / 2, 0]) {
    translate([0, 36, 0])
      cylinder(d = tab_screw_diameter, h = tab_thickness + 1);
    translate([0, -36, 0])
      cylinder(d = tab_screw_diameter, h = tab_thickness + 1);
  }
  translate([outer_length - tab_width / 2, outer_width / 2, 0]) {
    translate([0, 36, 0])
      cylinder(d = tab_screw_diameter, h = tab_thickness + 1);
    translate([0, -36, 0])
      cylinder(d = tab_screw_diameter, h = tab_thickness + 1);
  }
}

module vent_grid(l_num, w_num) {
  grid_spacing = grid_hole_width + grid_bar_width;
  width_offset = -(grid_spacing * w_num - grid_bar_width) / 2;
  length_offset = -(grid_spacing * l_num - grid_bar_width) / 2;
  translate([length_offset, width_offset, 0])
    for(len_pos = [0:grid_spacing:grid_spacing * (l_num - 1)])
      for(width_pos = [0:grid_spacing:grid_spacing * (w_num - 1)]) {
        translate([len_pos, width_pos, 0])
          cube([grid_hole_width, grid_hole_width, 50]);
      }
}

//translate([case_thickness + length_wall_gap, case_thickness + width_wall_gap, case_thickness + pcb_back_clearance])
//board();

//difference() {
//top();
//  translate([5,10.5,0])
//    cube([65,40,50]);
//
//}
//translate([0,0,bottom_height + 0.1])
//  top();

//rotate([0, 180, 0])
//translate([0,0,9])
top();
//bottom();

