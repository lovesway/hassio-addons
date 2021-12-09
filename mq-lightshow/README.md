# Home Assistant Add-on: MQ Light Show

## About
MQ-LightShow allows for the creation and control of light shows for lights 
which use Tasmota via the mqtt protocol. It was built to serve the needs of
[Love's Way](https://lovesway.info), but could support a broader scope
of needs if there is interest. It was frustrating to have all of these RGB 
lights around the ranch with no reasonable way to sync them all. Ultimately 
this app ended up being a great way to mitigate feral hog problems as well, 
this was done by creating a show that turns the outdoor lights on/off in 
sequence so as to create the illusion of a light moving back and forth around 
the ranch. Using mqtt does allow for a reasonable amount of speed and stability, 
things that are not possible using cloud based solutions.

Once the light shows are created, buttons can be added to the Home Assistant 
application for remote control. This also allows for automation. We use 
automations to activate an a ranch wide RGB fade scene at sundown, at 11pm 
the lights go into predator mode to detour the hogs/predators. On weekends 
they do a special party mode so we know it's Saturday night!

### Development Notes
This was not architected/designed and was developed out of necessity. While I did try to 
create abstractions for things that may need to be expanded, it is not in a pristine 
state. The ```.golangci.yml``` file contains hints as to what could be updated. 
Honeslty I didn't lint this code until a year after it was initially written. I 
didn't intend to publish it until my neighbor saw it and wanted to be able to use it. 
If it works for me I won't likely add features or updates unless it's really easy, I 
have a lot of chores and irons in the fire so to say. I'm absolutely open to contributors
though!

The UI does have some helpful features and I did seperate the API from the UI in the event 
that someone wanted to create a different UI at some point.

## License

MIT License

Copyright (c) 2020-2021 Brian McGowan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.