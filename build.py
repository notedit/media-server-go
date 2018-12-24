import os
import shutil



header_files = [
    "mediaserver/include/config.h",
    "mediaserver/include/dtls.h",
    "mediaserver/include/OpenSSL.h",
    "mediaserver/include/media.h",
    "mediaserver/include/rtp.h",
    "mediaserver/include/tools.h",
    "mediaserver/include/rtpsession.h",
    "mediaserver/include/DTLSICETransport.h",
    "mediaserver/include/RTPBundleTransport.h",
    "mediaserver/include/PCAPTransportEmulator.h",
    "mediaserver/include/mp4recorder.h",
    "mediaserver/include/mp4streamer.h",
    "mediaserver/src/vp9/VP9LayerSelector.h",
    "mediaserver/include/rtp/RTPStreamTransponder.h",
    "mediaserver/include/ActiveSpeakerDetector.h",
    "mediaserver/include/acumulator.h"
]


for src_file in header_files:
    dest_file = './include/' + src_file
    os.makedirs(os.path.dirname(dest_file),exist_ok=True)
    shutil.copy(src_file,dest_file)


