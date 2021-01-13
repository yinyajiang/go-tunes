#include "simpleApi.h"
#include "unkonwStruct.h"
#include "iTunesApi.h"



void AddLoadDir(wchar_t* dir)
{
#ifdef WIN32
	AddEnvLoadDir(dir);
#else

#endif
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//CoreFoundation.dll
#ifdef WIN32
void	CFRunLoopRun()
{
	CoreFoundation().CFRunLoopRun();
}

CFStringRef __CFStringMakeConstantString(const char *r1)
{
	return CoreFoundation().__CFStringMakeConstantString(r1);
}

void CFRelease(CFTypeRef r1)
{
	CoreFoundation().CFRelease(r1);
}

void*		kCFAllocatorDefault()
{
	return CoreFoundation().kCFAllocatorDefault;
}

void*		kCFTypeArrayCallBacks()
{
	return CoreFoundation().kCFTypeArrayCallBacks;
}

void*		kCFBooleanTrue()
{
	return CoreFoundation().kCFBooleanTrue;
}

void*		kCFBooleanFalse()
{
	return CoreFoundation().kCFBooleanFalse;
}

void*		kCFNumberPositiveInfinity()
{
	return CoreFoundation().kCFNumberPositiveInfinity;
}

void*		kCFTypeDictionaryKeyCallBacks()
{
	return CoreFoundation().kCFTypeDictionaryKeyCallBacks;
}

void*		kCFTypeDictionaryValueCallBacks()
{
	return CoreFoundation().kCFTypeDictionaryValueCallBacks;
}

CFMutableDictionaryRef CFDictionaryCreateMutable(int r1, int r2, void* r3, void* r4)
{
	return CoreFoundation().CFDictionaryCreateMutable(r1, r2, r3, r4);
}

CFMutableDictionaryRef CFDictionaryCreateMutableCopy(CFAllocatorRef r1, CFIndex r2, CFDictionaryRef r3)
{
	return CoreFoundation().CFDictionaryCreateMutableCopy(r1, r2, r3);
}

CFDictionaryRef		   CFDictionaryCreate(CFAllocatorRef r1, void* r2, void* r3, int r4, void* r5, void* r6)
{
	return CoreFoundation().CFDictionaryCreate(r1, r2, r3, r4, r5, r6);
}

CFDataRef CFPropertyListCreateXMLData(CFAllocatorRef r1, CFPropertyListRef r2)
{
	return CoreFoundation().CFPropertyListCreateXMLData(r1, r2);
}

int		 CFDataGetLength(CFDataRef r1)
{
	return CoreFoundation().CFDataGetLength(r1);
}

void*	CFDataGetBytePtr(CFDataRef r1)
{
	return CoreFoundation().CFDataGetBytePtr(r1);
}

CFMutableArrayRef CFArrayCreateMutable(CFAllocatorRef r1, int r2, void* r3)
{
	return CoreFoundation().CFArrayCreateMutable(r1, r2, r3);
}

CFMutableArrayRef CFArrayCreateMutableCopy(CFAllocatorRef r1, CFIndex r2, CFArrayRef r3)
{
	return CoreFoundation().CFArrayCreateMutableCopy(r1, r2, r3);
}

void CFArrayAppendArray(CFMutableArrayRef r1, CFArrayRef r2, CFRange r3)
{
	 CoreFoundation().CFArrayAppendArray(r1, r2, r3);
}

void CFDictionarySetValue(CFDictionaryRef r1, void* r2, void* r3)
{
	 CoreFoundation().CFDictionarySetValue(r1, r2, r3);
}

void CFDictionaryAddValue(CFDictionaryRef r1, void* r2, void* r3)
{
	 CoreFoundation().CFDictionaryAddValue(r1, r2, r3);
}

void* CFDictionaryGetValue(CFDictionaryRef r1, CFStringRef r2)
{
	return CoreFoundation().CFDictionaryGetValue(r1, r2);
}

int CFArrayGetCount(CFArrayRef r1)
{
	return CoreFoundation().CFArrayGetCount(r1);
}

void*  CFArrayGetValueAtIndex(CFArrayRef r1, int r2)
{
	return CoreFoundation().CFArrayGetValueAtIndex(r1, r2);
}

void CFArrayRemoveValueAtIndex(CFMutableArrayRef r1, int r2)
{
	 CoreFoundation().CFArrayRemoveValueAtIndex(r1, r2);
}

void CFArrayRemoveAllValues(CFMutableArrayRef r1)
{
	return CoreFoundation().CFArrayRemoveAllValues(r1);
}

CFNumberRef CFNumberCreate(CFAllocatorRef r1, CFNumberType r2, void* r3)
{
	return CoreFoundation().CFNumberCreate(r1, r2, r3);
}

void CFArrayAppendValue(CFMutableArrayRef r1, void* r2)
{
	return CoreFoundation().CFArrayAppendValue(r1, r2);
}

CFDataRef CFDataCreate(CFAllocatorRef r1, void* r2, int r3)
{
	return CoreFoundation().CFDataCreate(r1, r2, r3);
}

bool CFDictionaryContainsKey(CFDictionaryRef r1, void* r2)
{
	return CoreFoundation().CFDictionaryContainsKey(r1, r2);
}

CFStringRef CFStringCreateWithCString(CFAllocatorRef r1, const char * r2, int r3)
{
	return CoreFoundation().CFStringCreateWithCString(r1, r2, r3);
}

CFStringRef CFStringCreateWithCharacters(CFAllocatorRef r1, const wchar_t * r2, int r3)
{
	return CoreFoundation().CFStringCreateWithCharacters(r1, r2, r3);
}

CFStringRef CFStringCreateWithCharactersNoCopy(CFAllocatorRef r1, const wchar_t * r2, int r3, CFAllocatorRef r4)
{
	return CoreFoundation().CFStringCreateWithCharactersNoCopy(r1, r2, r3, r4);
}

int CFStringGetLength(CFStringRef r1)
{
	return CoreFoundation().CFStringGetLength(r1);
}
CFURLRef CFURLCreateWithFileSystemPath(CFAllocatorRef r1, CFStringRef r2, CFURLPathStyle r3, int r4)
{
	return CoreFoundation().CFURLCreateWithFileSystemPath(r1, r2, r3, r4);
}

CFReadStreamRef CFReadStreamCreateWithFile(CFAllocatorRef r1, CFURLRef r2)
{
	return CoreFoundation().CFReadStreamCreateWithFile(r1, r2);

}

int CFReadStreamOpen(CFReadStreamRef r1)
{
	return CoreFoundation().CFReadStreamOpen(r1);
}

void CFReadStreamClose(CFReadStreamRef r1)
{
	 CoreFoundation().CFReadStreamClose(r1);
}

CFPropertyListRef CFPropertyListCreateWithStream(CFAllocatorRef r1, CFReadStreamRef r2, CFIndex r3, int r4, void * r5, CFStringRef * r6)
{
	return CoreFoundation().CFPropertyListCreateWithStream(r1, r2, r3, r4, r5, r6);
}
CFMutableDataRef CFDataCreateMutable(CFAllocatorRef r1, CFIndex r2)
{
	return CoreFoundation().CFDataCreateMutable(r1, r2);
}
void CFDataAppendBytes(CFMutableDataRef r1, const uint8_t * r2, CFIndex r3)
{
	 CoreFoundation().CFDataAppendBytes(r1, r2, r3);
}
CFTypeID CFGetTypeID(CFTypeRef r1)
{
	return CoreFoundation().CFGetTypeID(r1);
}
CFPropertyListRef CFPropertyListCreateWithData(CFAllocatorRef r1, CFMutableDataRef r2 , int r3, void * r4, CFStringRef *r5)
{
	return CoreFoundation().CFPropertyListCreateWithData(r1, r2, r3, r4, r5);
}
CFPropertyListRef CFPropertyListCreateFromXMLData(CFAllocatorRef r1, CFDataRef r2, int r3, CFStringRef * r4)
{
	return CoreFoundation().CFPropertyListCreateFromXMLData(r1, r2, r3, r4);
}
int CFDictionaryGetValueIfPresent(CFDictionaryRef r1, void * r2, void ** r3)
{
	return CoreFoundation().CFDictionaryGetValueIfPresent(r1, r2, r3);
}
CFAbsoluteTime CFDateGetAbsoluteTime(CFDateRef r1)
{
	return CoreFoundation().CFDateGetAbsoluteTime(r1);
}
CFDateRef CFDateCreate(CFAllocatorRef r1, double r2)
{
	return CoreFoundation().CFDateCreate(r1, r2);
}
CFNumberType CFNumberGetType(CFNumberRef r1)
{
	return CoreFoundation().CFNumberGetType(r1);
}
int CFNumberGetValue(CFNumberRef r1, CFNumberType r2, void * r3)
{
	return CoreFoundation().CFNumberGetValue(r1, r2, r3);
}
int CFStringGetSystemEncoding()
{
	return CoreFoundation().CFStringGetSystemEncoding();
}
const char * CFStringGetCStringPtr(CFStringRef r1, CFStringEncoding r2)
{
	return CoreFoundation().CFStringGetCStringPtr(r1, r2);
}
int CFStringGetCString(CFStringRef r1, char * r2, int r3, CFStringEncoding r4)
{
	return CoreFoundation().CFStringGetCString(r1, r2, r3, r4);
}
int CFStringGetBytes(CFStringRef r1, CFRange r2, CFStringEncoding r3, uint8_t r4, bool r5, uint8_t * r6, CFIndex r7, CFIndex * r8)
{
	return CoreFoundation().CFStringGetBytes(r1, r2, r3, r4, r5, r6, r7, r8);
}
CFDataRef CFPropertyListCreateData(CFAllocatorRef r1, CFPropertyListRef r2, int r3, int r4, void ** r5)
{
	return CoreFoundation().CFPropertyListCreateData(r1, r2, r3, r4, r5);
}
int CFURLWriteDataAndPropertiesToResource(CFURLRef r1, CFDataRef r2, CFDictionaryRef r3, int * r4)
{
	return CoreFoundation().CFURLWriteDataAndPropertiesToResource(r1, r2, r3, r4);
}
int CFDictionaryGetCount(CFDictionaryRef r1)
{
	return CoreFoundation().CFDictionaryGetCount(r1);
}
void CFDictionaryGetKeysAndValues(CFDictionaryRef r1, void ** r2, void ** r3)
{
	 CoreFoundation().CFDictionaryGetKeysAndValues(r1, r2, r3);
}
CFTimeZoneRef CFTimeZoneCopyDefault()
{
	return CoreFoundation().CFTimeZoneCopyDefault();
}
CFTimeZoneRef CFTimeZoneCopySystem()
{
	return CoreFoundation().CFTimeZoneCopySystem();
}
CFTimeZoneRef CFTimeZoneCreateWithName(CFAllocatorRef r1, CFStringRef r2, bool r3)
{
	return CoreFoundation().CFTimeZoneCreateWithName(r1, r2, r3);
}
CFTimeZoneRef CFTimeZoneCreateWithTimeIntervalFromGMT(CFAllocatorRef r1, CFTimeInterval r2)
{
	return CoreFoundation().CFTimeZoneCreateWithTimeIntervalFromGMT(r1, r2);
}
CFAbsoluteTime CFAbsoluteTimeGetCurrent()
{
	return CoreFoundation().CFAbsoluteTimeGetCurrent();
}
CFGregorianDate CFAbsoluteTimeGetGregorianDate(CFAbsoluteTime r1, CFTimeZoneRef r2)
{
	return CoreFoundation().CFAbsoluteTimeGetGregorianDate(r1, r2);
}
CFAbsoluteTime CFGregorianDateGetAbsoluteTime(CFGregorianDate r1, CFTimeZoneRef r2)
{
	return CoreFoundation().CFGregorianDateGetAbsoluteTime(r1, r2);
}
CFTypeID CFStringGetTypeID()
{
	return CoreFoundation().CFStringGetTypeID();
}
CFTypeID CFDictionaryGetTypeID()
{
	return CoreFoundation().CFDictionaryGetTypeID();
}
CFTypeID CFDataGetTypeID()
{
	return CoreFoundation().CFDataGetTypeID();
}
CFTypeID CFNumberGetTypeID()
{
	return CoreFoundation().CFNumberGetTypeID();
}
CFTypeID CFAllocatorGetTypeID()
{
	return CoreFoundation().CFAllocatorGetTypeID();
}
CFTypeID CFURLGetTypeID()
{
	return CoreFoundation().CFURLGetTypeID();
}
CFTypeID CFReadStreamGetTypeID()
{
	return CoreFoundation().CFReadStreamGetTypeID();
}
void CFDictionaryReplaceValue(CFDictionaryRef r1, void* r2, void* r3)
{
	 CoreFoundation().CFDictionaryReplaceValue(r1, r2, r3);
}
CFTypeID CFArrayGetTypeID()
{
	return CoreFoundation().CFArrayGetTypeID();
}
CFTypeID CFDateGetTypeID()
{
	return CoreFoundation().CFDateGetTypeID();
}
CFTypeID CFErrorGetTypeID()
{
	return CoreFoundation().CFErrorGetTypeID();
}
CFTypeID CFNullGetTypeID()
{
	return CoreFoundation().CFNullGetTypeID();
}
CFTypeID CFBooleanGetTypeID()
{
	return CoreFoundation().CFBooleanGetTypeID();
}
CFTypeID CFAttributedStringGetTypeID()
{
	return CoreFoundation().CFAttributedStringGetTypeID();
}
CFTypeID CFBagGetTypeID()
{
	return CoreFoundation().CFBagGetTypeID();
}
CFTypeID CFBitVectorGetTypeID()
{
	return CoreFoundation().CFBitVectorGetTypeID();
}
CFTypeID CFBundleGetTypeID() 
{
	return CoreFoundation().CFBundleGetTypeID();
}
CFTypeID CFCalendarGetTypeID()
{
	return CoreFoundation().CFCalendarGetTypeID();
}
CFTypeID CFCharacterSetGetTypeID()
{
	return CoreFoundation().CFCharacterSetGetTypeID();
}
CFTypeID CFLocaleGetTypeID()
{
	return CoreFoundation().CFLocaleGetTypeID();
}
CFTypeID CFRunArrayGetTypeID()
{
	return CoreFoundation().CFRunArrayGetTypeID();
}
CFTypeID CFSetGetTypeID()
{
	return CoreFoundation().CFSetGetTypeID();
}
CFTypeID CFTimeZoneGetTypeID()
{
	return CoreFoundation().CFTimeZoneGetTypeID();
}
CFTypeID CFTreeGetTypeID()
{
	return CoreFoundation().CFTreeGetTypeID();
}
CFTypeID CFUUIDGetTypeID()
{
	return CoreFoundation().CFUUIDGetTypeID();
}
CFTypeID CFWriteStreamGetTypeID()
{
	return CoreFoundation().CFWriteStreamGetTypeID();
}
CFTypeID CFXMLNodeGetTypeID()
{
	return CoreFoundation().CFXMLNodeGetTypeID();
}
CFTypeID CFStorageGetTypeID() 
{
	return CoreFoundation().CFStorageGetTypeID();
}
CFTypeID CFSocketGetTypeID()
{
	return CoreFoundation().CFSocketGetTypeID();
}
CFTypeID CFWindowsNamedPipeGetTypeID()
{
	return CoreFoundation().CFWindowsNamedPipeGetTypeID();
}
CFTypeID CFPlugInGetTypeID()
{
	return CoreFoundation().CFPlugInGetTypeID();
}
CFTypeID CFPlugInInstanceGetTypeID()
{
	return CoreFoundation().CFPlugInInstanceGetTypeID();
}
CFTypeID CFBinaryHeapGetTypeID()
{
	return CoreFoundation().CFBinaryHeapGetTypeID();
}
CFTypeID CFDateFormatterGetTypeID()
{
	return CoreFoundation().CFDateFormatterGetTypeID();
}
CFTypeID CFMessagePortGetTypeID()
{
	return CoreFoundation().CFMessagePortGetTypeID();
}
CFTypeID CFNotificationCenterGetTypeID()
{
	return CoreFoundation().CFNotificationCenterGetTypeID();
}
CFTypeID CFNumberFormatterGetTypeID()
{
	return CoreFoundation().CFNumberFormatterGetTypeID();
}
CFTypeID _CFKeyedArchiverUIDGetTypeID()
{
	return CoreFoundation()._CFKeyedArchiverUIDGetTypeID();
}
int _CFKeyedArchiverUIDGetValue(void* r1)
{
	return CoreFoundation()._CFKeyedArchiverUIDGetValue(r1);
}
CFStringRef CFStringCreateWithFormat(CFAllocatorRef r1, CFDictionaryRef r2, CFStringRef r3, ...)
{
	va_list vl;
	va_start(vl, r3);
	auto ret = CoreFoundation().CFStringCreateWithFormat(r1, r2, r3,vl);
	va_end(vl);
	return ret;
}

CFBundleRef CFBundleGetMainBundle()
{
	return CoreFoundation().CFBundleGetMainBundle();
}
CFURLRef CFBundleCopyBundleURL(CFBundleRef r1)
{
	return CoreFoundation().CFBundleCopyBundleURL(r1);
}
CFURLRef CFURLCreateCopyDeletingLastPathComponent(CFAllocatorRef r1, CFURLRef r2)
{
	return CoreFoundation().CFURLCreateCopyDeletingLastPathComponent(r1, r2);
}
void* CFURLGetFileSystemRepresentation(CFURLRef r1, void* r2, uint8_t * r3, CFIndex r4)
{
	return CoreFoundation().CFURLGetFileSystemRepresentation(r1, r2, r3, r4);
}
#endif

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//MobileDevice.dll
int AMDeviceLookupApplications(void* r1, void* r2, void** r3)
{
	return  MobileDevice().AMDeviceLookupApplications(r1, r2, r3);
}
int AMDeviceStartHouseArrestService(void* r1, void* r2, void* r3, void* r4, void* r5)
{
	return  MobileDevice().AMDeviceStartHouseArrestService(r1, r2, r3, r4, r5);
}
int AMDeviceInstallApplication(void* r1, CFStringRef r2, void * r3, void * r4, void * r5)
{
	return  MobileDevice().AMDeviceInstallApplication(r1, r2, r3, r4, r5);
}
int AMDeviceUninstallApplication(void* r1, CFStringRef r2, void *r3, void *r4, void *r5)
{
	return  MobileDevice().AMDeviceUninstallApplication(r1, r2, r3, r4, r5);
}
int AMDeviceRemoveApplicationArchive(void* r1, CFStringRef r2, void * r3, void * r4, void *r5)
{
	return  MobileDevice().AMDeviceRemoveApplicationArchive(r1, r2, r3, r4, r5);
}
int AMDeviceArchiveApplication(void* r1, CFStringRef r2, void * r3, void * r4, void *r5)
{
	return  MobileDevice().AMDeviceArchiveApplication(r1, r2, r3, r4, r5);
}
int AFCConnectionOpen(void* r1, int r2, void** r3)
{
	return  MobileDevice().AFCConnectionOpen(r1, r2, r3);
}
int AMDServiceConnectionInvalidate(void* r1)
{
	return  MobileDevice().AMDServiceConnectionInvalidate(r1);
}
int AMDeviceNotificationSubscribe(void* r1, int r2, int r3, int r4, void** r5)
{
	return  MobileDevice().AMDeviceNotificationSubscribe(r1, r2, r3, r4, r5);
}
int AMDeviceNotificationUnsubscribe(void* r1)
{
	return  MobileDevice().AMDeviceNotificationUnsubscribe(r1);
}
int AMDeviceRelease(void* r1)
{
	return  MobileDevice().AMDeviceRelease(r1);
}
int AMDeviceConnect(void* r1)
{
	return  MobileDevice().AMDeviceConnect(r1);
}
int AMDeviceDisconnect(void* r1)
{
	return  MobileDevice().AMDeviceDisconnect(r1);
}
int AMDeviceIsPaired(void* r1)
{
	return MobileDevice().AMDeviceIsPaired(r1);
}
int AMDeviceValidatePairing(void* r1)
{
	return MobileDevice().AMDeviceValidatePairing(r1);
}
int AMDevicePair(void* r1)
{
	return MobileDevice().AMDevicePair(r1);
}
int AMDeviceUnpair(void* r1)
{
	return MobileDevice().AMDeviceUnpair(r1);
}
int AMDeviceStartSession(void* r1)
{
	return MobileDevice().AMDeviceStartSession(r1);
}
int AMDeviceSecureStartService(void* r1, void* r2, void* r3, void** r4)
{
	return MobileDevice().AMDeviceSecureStartService(r1, r2, r3, r4);
}
int AMDeviceStartService(void* r1, void* r2, void** r3, void* r4)
{
	return MobileDevice().AMDeviceStartService(r1, r2, r3, r4);
}
int AMDeviceStopSession(void* r1)
{
	return MobileDevice().AMDeviceStopSession(r1);
}
int AFCConnectionClose(void* r1)
{
	return MobileDevice().AFCConnectionClose(r1);
}
int AFCDeviceInfoOpen(void* r1, void** r2)
{
	return MobileDevice().AFCDeviceInfoOpen(r1, r2);
}
int AFCFileInfoOpen(void* r1, void* r2, void** r3)
{
	return MobileDevice().AFCFileInfoOpen(r1, r2, r3);
}
int AFCKeyValueRead(void* r1, void** r2, void** r3)
{
	return MobileDevice().AFCKeyValueRead(r1, r2, r3);
}
int AFCKeyValueClose(void* r1)
{
	return MobileDevice().AFCKeyValueClose(r1);
}
int AFCDirectoryOpen(void* r1, void* r2, void** r3)
{
	return MobileDevice().AFCDirectoryOpen(r1, r2, r3);
}
int AFCDirectoryRead(void* r1, void* r2, void** r3)
{
	return MobileDevice().AFCDirectoryRead(r1, r2, r3);
}
int AFCDirectoryClose(void* r1, void* r2)
{
	return MobileDevice().AFCDirectoryClose(r1, r2);
}
int AFCDirectoryCreate(void* r1, void* r2)
{
	return MobileDevice().AFCDirectoryCreate(r1, r2);
}
int AFCRemovePath(void* r1, void* r2)
{
	return MobileDevice().AFCRemovePath(r1, r2);
}
int AFCRenamePath(void* r1, void* r2, void* r3)
{
	return MobileDevice().AFCRenamePath(r1, r2, r3);
}
int AFCFileRefOpen(void* r1, void* r2, unsigned long long r3, unsigned long long* r4)
{
	return MobileDevice().AFCFileRefOpen(r1, r2, r3, r4);
}
int AFCFileRefRead(void* r1, unsigned long long r2, void* r3, void* r4)
{
	return MobileDevice().AFCFileRefRead(r1, r2, r3, r4);
}
int AFCFileRefWrite(void* r1, unsigned long long r2, void* r3, int r4)
{
	return MobileDevice().AFCFileRefWrite(r1, r2, r3, r4);
}
int AFCFileRefClose(void* r1, unsigned long long r2)
{
	return MobileDevice().AFCFileRefClose(r1, r2);
}
int AFCFileRefSeek(void* r1, unsigned long long r2, unsigned long long r3, unsigned long r4)
{
	return MobileDevice().AFCFileRefSeek(r1, r2, r3, r4);
}
int AFCFileRefTell(void* r1, unsigned long long r2, unsigned long * r3)
{
	return MobileDevice().AFCFileRefTell(r1, r2, r3);
}
void* AMDeviceCopyDeviceIdentifier(void* r1)
{
	return MobileDevice().AMDeviceCopyDeviceIdentifier(r1);
}
void* AMDeviceCopyValue(void* r1, void* r2, void* r3)
{
	return MobileDevice().AMDeviceCopyValue(r1, r2, r3);
}
int AMDeviceGetInterfaceType(void* r1)
{
	return MobileDevice().AMDeviceGetInterfaceType(r1);
}
int AMRestoreRegisterForDeviceNotifications(void* r1, void* r2, void* r3, void* r4, int r5, void* r6)
{
	return MobileDevice().AMRestoreRegisterForDeviceNotifications(r1, r2, r3, r4, r5, r6);
}
int USBMuxConnectByPort(int r1, short r2, void* r3)
{
	return MobileDevice().USBMuxConnectByPort(r1, r2, r3);
}
int AMRestorePerformRecoveryModeRestore(void* r1, void* r2, void* r3, void* r4)
{
	return MobileDevice().AMRestorePerformRecoveryModeRestore(r1, r2, r3, r4);
}
int AMRestorePerformDFURestore(void* r1, void* r2, void* r3, void* r4)
{
	return MobileDevice().AMRestorePerformDFURestore(r1, r2, r3, r4);
}
int AMRestorableDeviceRegisterForNotificationsForDevices(am_recovery_device_notification_callback r1, void* r2, unsigned int r3, void* r4, void* r5)
{
	return MobileDevice().AMRestorableDeviceRegisterForNotificationsForDevices(r1, r2, r3, r4, r5);
}
void AMRestoreUnregisterForDeviceNotifications()
{
	 MobileDevice().AMRestoreUnregisterForDeviceNotifications();
}
int AMRestorableDeviceRestore(am_restore_device* r1, CFDictionaryRef r2, void* r3, void* r4)
{
	return MobileDevice().AMRestorableDeviceRestore(r1, r2, r3, r4);
}
int AMSRestoreWithApplications(void* r1, void* r2, void* r3, void* r4, void* r5, void* r6, void* r7)
{
	return MobileDevice().AMSRestoreWithApplications(r1, r2, r3, r4, r5, r6, r7);
}
int AMSUnregisterTarget(void* r1)
{
	return MobileDevice().AMSUnregisterTarget(r1);
}
int AMDeviceSetValue(void* r1, void* r2, void* r3, void* r4)
{
	return MobileDevice().AMDeviceSetValue(r1, r2, r3, r4);
}
int AMRecoveryModeDeviceSendFileToDevice(void* r1, CFStringRef r2)
{
	return MobileDevice().AMRecoveryModeDeviceSendFileToDevice(r1, r2);
}
int AMRecoveryModeDeviceSendCommandToDevice(void* r1, CFStringRef r2)
{
	return MobileDevice().AMRecoveryModeDeviceSendCommandToDevice(r1, r2);
}
uint16_t AMRecoveryModeDeviceGetProductID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetProductID(r1);
}
unsigned long AMRecoveryModeDeviceGetProductType(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetProductType(r1);
}
unsigned long AMRecoveryModeDeviceGetChipID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetChipID(r1);
}
uint64_t AMRecoveryModeDeviceGetECID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetECID(r1);
}
unsigned long AMRecoveryModeDeviceGetLocationID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetLocationID(r1);
}
unsigned long AMRecoveryModeDeviceGetBoardID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetBoardID(r1);
}
unsigned char AMRecoveryModeDeviceGetProductionMode(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetProductionMode(r1);
}
unsigned long AMRecoveryModeDeviceGetTypeID(void* r1)
{
	return MobileDevice().AMRecoveryModeDeviceGetTypeID(r1);
}
void* AMRecoveryModeGetSoftwareBuildVersion(void* r1)
{
	return MobileDevice().AMRecoveryModeGetSoftwareBuildVersion(r1);
}
uint16_t AMDFUModeDeviceGetProductID(void* r1)
{
	return MobileDevice().AMDFUModeDeviceGetProductID(r1);
}
unsigned long AMDFUModeDeviceGetProductType(void* r1)
{
	return MobileDevice().AMDFUModeDeviceGetProductType(r1);
}
unsigned long AMDFUModeDeviceGetChipID(void* r1)
{
	return MobileDevice().AMDFUModeDeviceGetChipID(r1);
}
uint64_t AMDFUModeDeviceGetECID(void* r1){

	return MobileDevice().AMDFUModeDeviceGetECID(r1);
}
unsigned long AMDFUModeDeviceGetLocationID(void* r1){

	return MobileDevice().AMDFUModeDeviceGetLocationID(r1);
}
unsigned long AMDFUModeDeviceGetBoardID(void* r1){

	return MobileDevice().AMDFUModeDeviceGetBoardID(r1);
}
unsigned char AMDFUModeDeviceGetProductionMode(void* r1){

	return MobileDevice().AMDFUModeDeviceGetProductionMode(r1);
}
unsigned long AMDFUModeDeviceGetTypeID(void* r1){

	return MobileDevice().AMDFUModeDeviceGetTypeID(r1);
}
int AMRecoveryModeDeviceSetAutoBoot(void* r1, unsigned char r2){

	return MobileDevice().AMRecoveryModeDeviceSetAutoBoot(r1, r2);
}
int AMRecoveryModeDeviceReboot(void* r1){

	return MobileDevice().AMRecoveryModeDeviceReboot(r1);
}
int AMRestoreModeDeviceReboot(void* r1){

	return MobileDevice().AMRestoreModeDeviceReboot(r1);
}
int AMRestoreEnableFileLogging(char* r1){

	return MobileDevice().AMRestoreEnableFileLogging(r1);
}
int AMRestoreDisableFileLogging(){

	return MobileDevice().AMRestoreDisableFileLogging();
}
int AMRestorableDeviceGetState(void* r1){

	return MobileDevice().AMRestorableDeviceGetState(r1);
}
void* AMRestorableDeviceCopyDFUModeDevice(void* r1){

	return MobileDevice().AMRestorableDeviceCopyDFUModeDevice(r1);
}
void* AMRestorableDeviceCopyRecoveryModeDevice(void* r1){

	return MobileDevice().AMRestorableDeviceCopyRecoveryModeDevice(r1);
}
void* AMRestorableDeviceCopyAMDevice(void* r1){

	return MobileDevice().AMRestorableDeviceCopyAMDevice(r1);
}
void* AMRestorableDeviceCreateFromAMDevice(void* r1){

	return MobileDevice().AMRestorableDeviceCreateFromAMDevice(r1);
}
uint16_t AMRestorableDeviceGetProductID(void* r1){

	return MobileDevice().AMRestorableDeviceGetProductID(r1);
}
unsigned long AMRestorableDeviceGetProductType(void* r1){

	return MobileDevice().AMRestorableDeviceGetProductType(r1);
}
unsigned long AMRestorableDeviceGetChipID(void* r1){

	return MobileDevice().AMRestorableDeviceGetChipID(r1);
}
uint64_t AMRestorableDeviceGetECID(void* r1){

	return MobileDevice().AMRestorableDeviceGetECID(r1);
}
unsigned long AMRestorableDeviceGetLocationID(void* r1){

	return MobileDevice().AMRestorableDeviceGetLocationID(r1);
}
unsigned long AMRestorableDeviceGetBoardID(void* r1){

	return MobileDevice().AMRestorableDeviceGetBoardID(r1);
}
unsigned long AMRestoreModeDeviceGetTypeID(void* r1){

	return MobileDevice().AMRestoreModeDeviceGetTypeID(r1);
}
void* AMRestoreModeDeviceCopySerialNumber(void* r1){

	return MobileDevice().AMRestoreModeDeviceCopySerialNumber(r1);
}
void* AMRestorableDeviceCopySerialNumber(void* r1){

	return MobileDevice().AMRestorableDeviceCopySerialNumber(r1);
}
void* AMRecoveryModeDeviceCopySerialNumber(void* r1){

	return MobileDevice().AMRecoveryModeDeviceCopySerialNumber(r1);
}
int AFCConnectionGetContext(void* r1){

	return MobileDevice().AFCConnectionGetContext(r1);
}
int AFCConnectionGetFSBlockSize(void* r1){

	return MobileDevice().AFCConnectionGetFSBlockSize(r1);
}
int AFCConnectionGetIOTimeout(void* r1){

	return MobileDevice().AFCConnectionGetIOTimeout(r1);
}
int AFCConnectionGetSocketBlockSize(void* r1){

	return MobileDevice().AFCConnectionGetSocketBlockSize(r1);
}
void* AMRestoreCreateDefaultOptions(void* r1){

	return MobileDevice().AMRestoreCreateDefaultOptions(r1);
}
int AMRestorePerformRestoreModeRestore(void* r1, void* r2, void* r3, void* r4){

	return MobileDevice().AMRestorePerformRestoreModeRestore(r1, r2, r3, r4);
}
void* AMRestoreModeDeviceCreate(int r1, int r2, int r3)
{
	return MobileDevice().AMRestoreModeDeviceCreate(r1, r2, r3);
}
int AMRestoreCreatePathsForBundle(void* r1, void* r2, void* r3, int r4, void** r5, void** r6, int r7, void** r8)
{
	return MobileDevice().AMRestoreCreatePathsForBundle(r1, r2, r3, r4, r5, r6, r7, r8);
}
int AMDeviceGetConnectionID(void* r1)
{
	return MobileDevice().AMDeviceGetConnectionID(r1);
}
int AMDeviceEnterRecovery(void* r1)
{
	return MobileDevice().AMDeviceEnterRecovery(r1);
}
int AMDeviceRetain(void* r1)
{
	return MobileDevice().AMDeviceRetain(r1);
}
int AMDShutdownNotificationProxy(void* r1)
{
	return MobileDevice().AMDShutdownNotificationProxy(r1);
}
int AMDeviceDeactivate(void* r1)
{
	return MobileDevice().AMDeviceDeactivate(r1);
}
int AMDeviceActivate(void* r1, void* r2)
{
	return MobileDevice().AMDeviceActivate(r1, r2);
}
int AMDeviceRemoveValue(void* r1, int r2, void* r3)
{
	return MobileDevice().AMDeviceRemoveValue(r1, r2, r3);
}
int USBMuxListenerCreate(void* r1, void** r2)
{
	return MobileDevice().USBMuxListenerCreate(r1, r2);
}
int USBMuxListenerHandleData(void* r1)
{
	return MobileDevice().USBMuxListenerHandleData(r1);
}
int AMDObserveNotification(void* r1, void* r2)
{
	return MobileDevice().AMDObserveNotification(r1, r2);
}
int AMSInitialize()
{
	return MobileDevice().AMSInitialize();
}
int AMDListenForNotifications(void* r1, void* r2, void* r3)
{
	return MobileDevice().AMDListenForNotifications(r1, r2, r3);
}
int AMDeviceStartServiceWithOptions(void* r1, CFStringRef r2, void* r3, int* r4)
{
	return MobileDevice().AMDeviceStartServiceWithOptions(r1, r2, r3, r4);
}
void* AMDServiceConnectionCreate(CFStringRef r1, void* r2, void* r3)
{
	return MobileDevice().AMDServiceConnectionCreate(r1, r2, r3);
}
void* AMDServiceConnectionGetSocket(void* r1)
{
	return MobileDevice().AMDServiceConnectionGetSocket(r1);
}
int AMDServiceConnectionGetSecureIOContext(void* r1)
{
	return MobileDevice().AMDServiceConnectionGetSecureIOContext(r1);
}
int AMDServiceConnectionReceive(void* r1, void* r2, int r3)
{
	return MobileDevice().AMDServiceConnectionReceive(r1, r2, r3);
}
int AMDServiceConnectionSend(void* r1, void* r2, int r3)
{
	return MobileDevice().AMDServiceConnectionSend(r1, r2, r3);
}
int AMDServiceConnectionReceiveMessage(void* r1, CFDictionaryRef * r2, int * r3)
{
	return MobileDevice().AMDServiceConnectionReceiveMessage(r1, r2, r3);
}
int AMDServiceConnectionSendMessage(void* r1, CFDictionaryRef r2, int r3)
{
	return MobileDevice().AMDServiceConnectionSendMessage(r1, r2, r3);
}
int AMSChangeBackupPassword(CFStringRef r1, CFStringRef r2, CFStringRef r3, int * r4)
{
	return MobileDevice().AMSChangeBackupPassword(r1, r2, r3, r4);
}
int AMSBackupWithOptions(CFStringRef r1, CFStringRef r2, CFStringRef r3, CFDictionaryRef r4, BackupCallBack r5, int r6)
{
	return MobileDevice().AMSBackupWithOptions(r1, r2, r3, r4, r5, r6);
}
int AMSCancelBackupRestore(void * r1)
{
	return MobileDevice().AMSCancelBackupRestore(r1);
}
CFStringRef AMSGetErrorReasonForErrorCode(int r1)
{
	return MobileDevice().AMSGetErrorReasonForErrorCode(r1);
}


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//AirTrafficHost.dll
void* ATCFMessageGetParam(void* r1, void* r2)
{
	return AirTrafficHost().ATCFMessageGetParam(r1, r2);
}
int ATHostConnectionGetCurrentSessionNumber(void* r1) 
{
	return AirTrafficHost().ATHostConnectionGetCurrentSessionNumber(r1);
}
void ATHostConnectionSendFileProgress(void* r1, void* r2, void* r3, double r4, int r5, int r6)
{
	 AirTrafficHost().ATHostConnectionSendFileProgress(r1, r2, r3, r4, r5, r6);
}
void* ATCFMessageCreate(int r1, void* r2, void* r3)
{
	return AirTrafficHost().ATCFMessageCreate(r1, r2, r3);
}
void* ATHostConnectionCreateWithLibrary(void* r1, void* r2, void* r3)
{
	return AirTrafficHost().ATHostConnectionCreateWithLibrary(r1, r2, r3);
}
int ATHostConnectionCreate(void* r1)
{
	return AirTrafficHost().ATHostConnectionCreate(r1);
}
void ATHostConnectionSendPing(void* r1)
{
	 AirTrafficHost().ATHostConnectionSendPing(r1);
}
void ATHostConnectionSendAssetMetricsRequest(void* r1, int r2)
{
	 AirTrafficHost().ATHostConnectionSendAssetMetricsRequest(r1, r2);
}
void ATHostConnectionInvalidate(void* r1)
{
	 AirTrafficHost().ATHostConnectionInvalidate(r1);
}
void ATHostConnectionClose(void* r1)
{
	 AirTrafficHost().ATHostConnectionClose(r1);
}
void ATHostConnectionRelease(void* r1)
{
	AirTrafficHost().ATHostConnectionRelease(r1);
}
int ATHostConnectionSendPowerAssertion(void* r1, void* r2)
{
	return AirTrafficHost().ATHostConnectionSendPowerAssertion(r1, r2);
}
int ATHostConnectionRetain(void* r1)
{
	return AirTrafficHost().ATHostConnectionRetain(r1);
}
int ATHostConnectionSendMetadataSyncFinished(void* r1, void* r2, void* r3)
{
	return AirTrafficHost().ATHostConnectionSendMetadataSyncFinished(r1, r2, r3);
}
void ATHostConnectionSendFileError(void* r1, void* r2, void* r3, int r4)
{
	AirTrafficHost().ATHostConnectionSendFileError(r1, r2, r3, r4);
}
int ATHostConnectionSendAssetCompleted(void* r1, void* r2, void* r3, void* r4)
{
	return AirTrafficHost().ATHostConnectionSendAssetCompleted(r1, r2, r3, r4);
}
void* ATCFMessageGetName(void* r1)
{
	return AirTrafficHost().ATCFMessageGetName(r1);
}
int ATHostConnectionSendHostInfo(void* r1, void* r2)
{
	return AirTrafficHost().ATHostConnectionSendHostInfo(r1, r2);
}
int ATHostConnectionSendSyncRequest(void* r1, void* r2, void* r3, void* r4)
{
	return AirTrafficHost().ATHostConnectionSendSyncRequest(r1, r2, r3, r4);
}
int ATHostConnectionSendMessage(void* r1, void* r2)
{
	return AirTrafficHost().ATHostConnectionSendMessage(r1, r2);
}
int ATHostConnectionGetGrappaSessionId(int r1)
{
	return AirTrafficHost().ATHostConnectionGetGrappaSessionId(r1);
}
void* ATHostConnectionReadMessage(void* r1)
{
	return AirTrafficHost().ATHostConnectionReadMessage(r1);
}
