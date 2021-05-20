# baetyl-video-infer

The baetyl-video-infer module is an official module of [BAETYL](https://baetyl.io), which is used for capturing video frame and AI model inference. For capturing video frame, baetyl-video-infer module can support IP-network camera, USB camera and video files. And the supported DL(Deep Learning) frameworks and layers are keeping in sync with OpenCV, more detailed can be found [here](https://github.com/opencv/opencv/wiki/Deep-Learning-in-OpenCV).

In addition, baetyl-video-infer module is highly integrated with other modules of BAETYL framework empowering the edge AI and some industry scenes. Besides, baetyl-video-infer also provides accelerated processing support for AI models based on CPU and specific hardware (such as OpenCL(OpenVINO, INFERENCE_ENGINE)).

## Build Program

The baetyl-video-infer module is depend on [GoCV](https://github.com/hybridgroup/gocv). Please make sure GoCV works properly before build.

```shell
git clone https://github.com/baetyl/baetyl-video-infer.git
cd baetyl-video-infer
make build # build baetyl-video-infer
```

## Build Image

As above description, baetyl-video-infer module now can support CPU and specific hardware(Intel GPU, as known as OpenCL) for AI models inference. For the CPU, baetyl can support Linux-armv7, Linux-arm64, Linux-amd64 platforms, and provides accelerated processing support for AI models due to Intel GPU(OpenCL) of [OpenVINO framework](https://docs.openvinotoolkit.org/latest/index.html).

```shell
# CPU
make image # build image for CPU

# OpenVINO(OpenCL)
make image-openvino # build image for OpenVINO

# output images 
docker images | grep baetyl/video-infer
# baetyl/video-infer-openvino                                      git-e7ed6d0                                      4d89e690b226        About an hour ago   2.53GB
# baetyl/video-infer                                               git-e7ed6d0                                      ccab5f4b07ff        3 hours ago         672MB
```

**NOTE**: 

- For OpenVINO(OpenCL), now only support linux/amd64 platform;
- For CPU, cross compile is not support.

## Configuration

```yaml
video:
  uri: [MUST] The video file path or camera address. 
    # For IP camera, the configuration just like `rtsp://<username>:<password>@<ip>:<port>/Streaming/channels/<stream_number>/`
      # `<username>` and `<password>` are the login authentication element
      # `<ip>` is the IP-address of camera
      # `<port>` is the port number of RTSP protocol, the default value is `554`
      # `<stream_number>` is the channel number, if it is equal to `1`, it indicates that the main stream is being captured; if it is equal to `2`, it indicates that the secondary stream is being captured
    # For USB camera, the configuration just like "0"(represents mapping device `/dev/video0` into container, also should be mounted on video infer service)
    # For video file, the configuration just like `var/db/baetyl/data/test.mp4`(mount the volume(store the video file) on video infer service)
  limit:
    fps: [MUST] The max number of video frames handled by inference per second. If the video fps is N, limit.fps is M, then Ceil(N/M) - 1 frames will be skipped.
process: 
  before: creates 4-dimensional blob from image. Optionally resizes and crops image from center, subtract mean values, scales values by scalefactor, swap Blue and Red channels. More detailed contents please refer to https://docs.opencv.org/4.1.1/d6/d0f/group__dnn.html#ga29f34df9376379a603acd8df581ac8d7.
    scale: multiplier for image values.
    swaprb: flag which indicates that swap first and last channels in 3-channel image is necessary. 
    width: width of spatial size for output image.
    height: height of spatial size for output image.
    mean: scalar with mean values which are subtracted from channels. Values are intended to be in (mean-R, mean-G, mean-B) order if image has BGR ordering and swapRB is true.
      v1: blue component of type Scalar(Scalar is a 4-element(v1, v2, v3, v4) vector widely used in OpenCV to pass pixel values).
      v2: green component of type Scalar.
      v3: red component of type Scalar.
      v4: alpha component of type Scalar.
    crop: flag which indicates whether image will be cropped after resize or not.
  infer:
    model: [MUST] The path of model file, more detailed contents please refer to https://docs.opencv.org/4.1.1/d6/d0f/group__dnn.html#ga3b34fe7a29494a6a4295c169a7d32422.
    config: [MUST] The path of model config file, more detailed contents please refer to https://docs.opencv.org/4.1.1/d6/d0f/group__dnn.html#ga3b34fe7a29494a6a4295c169a7d32422.
    backend: [Optional] The network backend which is used to improve inference efficiency. Now support `halide`, `openvino`, `opencv`, `vulkan` and `default`. More detailed contents please refer to https://docs.opencv.org/4.1.1/d6/d0f/group__dnn.html#ga186f7d9bfacac8b0ff2e26e2eab02625.
    device: [Optional] The target device of DNN processing. Now support `cpu`(default), `fp32`, `fp16`, `vpu`, `vulkan` and `fpga`. More detailed contents please refer to https://docs.opencv.org/4.1.1/d6/d0f/group__dnn.html#ga709af7692ba29788182cf573531b0ff5.
  after:
    function: 
      name: [MUST] The name of the function that handle the inference result.
logger:
  filename: If the path is specified, writes log to the file, otherwise writes to stdout.
  level: The log level, support `debug`、`info`(default)、`warn` and `error`.
