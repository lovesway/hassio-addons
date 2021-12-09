## Examples

### Adding LightShows to Home Assistant
This is example code that would go into the configuration.yaml

```
switch:
  - platform: mqtt
    name: "Light Show Full Ranch Fade"
    state_topic: "mqlightshow/show/full_ranch_fade/stat"
    command_topic: "mqlightshow/show/full_ranch_fade/cmnd"
    payload_on: "ON"
    payload_off: "OFF"
    state_on: "ON"
    state_off: "OFF"
    optimistic: false
    qos: 0
    retain: true
  - platform: mqtt
    name: "Light Show Predator"
    state_topic: "mqlightshow/show/predator/stat"
    command_topic: "mqlightshow/show/predator/cmnd"
    payload_on: "ON"
    payload_off: "OFF"
    state_on: "ON"
    state_off: "OFF"
    optimistic: false
    qos: 0
    retain: true
```

The Entities card configuration example would look like this.
```
entities:
  - entity: switch.light_show_full_ranch_fade
  - entity: switch.light_show_predator
show_header_toggle: false
title: Light Shows
type: entities
```

### Known Limitations
Currently there is only support for Tasmota commands, but the device types are 
abstracted so that other firmware types/commands could be added.

MQTT configuration only supports ```%topic%/%prefix%/``` vs ```%prefix%/%topic%/```. 
This could be changed, but at this point it is a known limitation.

There is no support for importing/exporting or backing up the database. This is partially 
due to each isntallation having different devices, but I'm sure there are solutions 
for all of that. What is possible, is backing up ther sqlite.db manually or even accessing
it via cli sqlite. This file is likely located somewhere in /usr/share/hassio/addons/data.
