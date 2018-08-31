#Installation instructions for Raspbian

Install OpenCV+gocv using the documentation found at https://github.com/hybridgroup/gocv/blob/master/README.md (do not install the regular OpenCV packages).

On the Raspberry Pi, ensure the v4l2 driver is loaded.
```
modprobe bcm2835-v4l2
```