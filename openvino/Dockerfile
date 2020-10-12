FROM hub.baidubce.com/baetyl/openvino-go-devel:2020.4 as devel
COPY * /root/go/baetyl-video-infer/
RUN cd /root/go/baetyl-video-infer/ && \
    make build GO_TAGS="-tags openvino"

FROM hub.baidubce.com/baetyl/openvino-runtime:2020.4

USER root
WORKDIR /

# baetyl-video-infer
COPY --from=devel /root/go/baetyl-video-infer/output/baetyl-video-infer /bin/

ENV INTEL_OPENVINO_DIR "/opt/intel/openvino"
ENV INTEL_CVSDK_DIR "${INTEL_OPENVINO_DIR}"
ENV HDDL_INSTALL_DIR "${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/hddl"
ENV LD_LIBRARY_PATH "${INTEL_OPENVINO_DIR}/opencv/lib:${INTEL_OPENVINO_DIR}/deployment_tools/ngraph/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/hddl_unite/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/hddl/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/gna/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/mkltiny_lnx/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/external/tbb/lib:${INTEL_OPENVINO_DIR}/deployment_tools/inference_engine/lib/intel64"
ENV LIBRARY_PATH "${INTEL_OPENVINO_DIR}/data_processing/dl_streamer/lib:${INTEL_OPENVINO_DIR}/data_processing/gstreamer/lib"
ENV PATH "${INTEL_OPENVINO_DIR}/deployment_tools/model_optimizer:${INTEL_OPENVINO_DIR}/data_processing/gstreamer/bin:${PATH}"
ENV OpenCV_DIR "${INTEL_CVSDK_DIR}/opencv/cmake"

ENV CGO_CXXFLAGS "--std=c++11"
ENV CGO_CPPFLAGS "-I${INTEL_CVSDK_DIR}/opencv/include -I${INTEL_CVSDK_DIR}/deployment_tools/inference_engine/include"
ENV CGO_LDFLAGS "-L${INTEL_CVSDK_DIR}/opencv/lib -L${INTEL_CVSDK_DIR}/deployment_tools/inference_engine/lib/intel64 -lpthread -ldl -lHeteroPlugin -lMKLDNNPlugin -lmyriadPlugin -linference_engine -lclDNNPlugin -lopencv_core -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_calib3d -lopencv_flann -lopencv_gapi -lopencv_ml -lopencv_photo -lopencv_stitching -lopencv_videoio_ffmpeg -lopencv_videoio_gstreamer -lopencv_videoio_intel_mfx"

ENTRYPOINT ["baetyl-video-infer"]