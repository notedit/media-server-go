
class PlayerListener {
public:
	virtual ~PlayerListener() {}
	virtual void onEnd() {}
};

class REMBListener {
public:
	virtual ~REMBListener() {}
	virtual void onREMB() {}
};

class TargetBitrateListener {
public:
	virtual ~TargetBitrateListener() {}
	virtual void onBitrate() {}
};



