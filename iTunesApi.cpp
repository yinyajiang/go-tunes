#include "iTunesApi.h"
#include <string.h>
#include <vector>


void AddEnvLoadDir(std::wstring dir)
{
#ifdef WIN32
	DWORD dwSizeNeeded = ::GetEnvironmentVariable(L"PATH", NULL, 0);
	dwSizeNeeded += dir.length() + 8; 
	std::vector<wchar_t> buff;
	buff.resize(dwSizeNeeded);
	memset(&buff[0], 0, buff.size() * sizeof(wchar_t));
	::GetEnvironmentVariable(L"PATH", (LPTSTR)&buff[0], dwSizeNeeded);

	std::wstring env = (LPTSTR)&buff[0];
	env += L";";
	env += dir;
	::SetEnvironmentVariable(L"PATH", dir.c_str());
#endif
}


CCoreFoundation& CoreFoundation()
{
	static CCoreFoundation obj;
	return obj;
}

CMobileDevice&	MobileDevice(bool newdll)
{
	static CMobileDevice obj(newdll);
	return obj;
}

CAirTrafficHost&    AirTrafficHost()
{
	static CAirTrafficHost obj;
	return obj;
}


#ifdef WIN32
#define FUNC_LOAD(_,func)		   this->func = (PF_##func)GetProcAddress(m_hDll,#func);
#define VALUE_LOAD(func)		   FUNC_LOAD(0,func);\
									{ \
										void** ppTmp = (void**)this->func; \
										if (ppTmp) \
											this->func = (PF_##func)*ppTmp; \
									}
					
#else
#define FUNC_LOAD(direct,func) \
				if (direct)\ 
				{ \
					this->func = (PF_##func)::func; \
				} \
				else {\
					this->func = (PF_##func)::CFBundleGetFunctionPointerForName(pModle, CFStringCreateWithCString(NULL, #func, kCFStringEncodingUTF8));\
				}
#define VALUE_LOAD(func)		   this->func = (PF_##func)(&::func);
#endif



CCoreFoundation::CCoreFoundation(void)
{
#ifdef WIN32
	m_hDll = LoadLibrary(L"CoreFoundation.dll");
	if (NULL == m_hDll)
		MessageBoxA(NULL, "CoreFoundation.dll load fail", "error", 0);
#endif
	VALUE_LOAD(kCFAllocatorDefault);
	VALUE_LOAD(kCFTypeArrayCallBacks);
	VALUE_LOAD(kCFBooleanTrue);
	VALUE_LOAD(kCFBooleanFalse);
	VALUE_LOAD(kCFNumberPositiveInfinity);
	VALUE_LOAD(kCFTypeDictionaryKeyCallBacks);
	VALUE_LOAD(kCFTypeDictionaryValueCallBacks);

	FUNC_LOAD(1, CFRunLoopRun);
	FUNC_LOAD(1, __CFStringMakeConstantString);
	FUNC_LOAD(1, CFRelease);
	FUNC_LOAD(1, CFDictionaryCreateMutable);
	FUNC_LOAD(1, CFDictionaryCreateMutableCopy);
	FUNC_LOAD(1, CFDictionaryCreate);
	FUNC_LOAD(1, CFPropertyListCreateXMLData);
	FUNC_LOAD(1, CFDataGetLength);
	FUNC_LOAD(1, CFDataGetBytePtr);
	FUNC_LOAD(1, CFArrayCreateMutable);
	FUNC_LOAD(1, CFArrayCreateMutableCopy);
	FUNC_LOAD(1, CFArrayAppendArray);
	FUNC_LOAD(1, CFDictionarySetValue);
	FUNC_LOAD(1, CFDictionaryAddValue);
	FUNC_LOAD(1, CFDictionaryGetValue);
	FUNC_LOAD(1, CFArrayGetCount);
	FUNC_LOAD(1, CFArrayGetValueAtIndex);
	FUNC_LOAD(1, CFArrayRemoveValueAtIndex);
	FUNC_LOAD(1, CFArrayRemoveAllValues);
	FUNC_LOAD(1, CFNumberCreate);
	FUNC_LOAD(1, CFArrayAppendValue);
	FUNC_LOAD(1, CFDataCreate);
	FUNC_LOAD(1, CFDictionaryContainsKey);
	FUNC_LOAD(1, CFStringCreateWithCString);
	FUNC_LOAD(1, CFStringCreateWithCharacters);
	FUNC_LOAD(1, CFStringCreateWithCharactersNoCopy);
	FUNC_LOAD(1, CFStringGetLength);
	FUNC_LOAD(1, CFURLCreateWithFileSystemPath);
	FUNC_LOAD(1, CFReadStreamCreateWithFile);
	FUNC_LOAD(1, CFReadStreamOpen);
	FUNC_LOAD(1, CFReadStreamClose);
	FUNC_LOAD(1, CFPropertyListCreateWithStream);
	FUNC_LOAD(1, CFDataCreateMutable);
	FUNC_LOAD(1, CFDataAppendBytes);
	FUNC_LOAD(1, CFGetTypeID);
	FUNC_LOAD(1, CFPropertyListCreateWithData);
	FUNC_LOAD(1, CFPropertyListCreateFromXMLData);
	FUNC_LOAD(1, CFDictionaryGetValueIfPresent);
	FUNC_LOAD(1, CFDateGetAbsoluteTime);
	FUNC_LOAD(1, CFDateCreate);
	FUNC_LOAD(1, CFNumberGetType);
	FUNC_LOAD(1, CFNumberGetValue);
	FUNC_LOAD(1, CFStringGetSystemEncoding);
	FUNC_LOAD(1, CFStringGetCStringPtr);
	FUNC_LOAD(1, CFStringGetCString);
	FUNC_LOAD(1, CFStringGetBytes);
	FUNC_LOAD(1, CFPropertyListCreateData);
	FUNC_LOAD(1, CFURLWriteDataAndPropertiesToResource);
	FUNC_LOAD(1, CFDictionaryGetCount);
	FUNC_LOAD(1, CFDictionaryGetKeysAndValues);
	FUNC_LOAD(1, CFTimeZoneCopyDefault);
	FUNC_LOAD(1, CFTimeZoneCopySystem);
	FUNC_LOAD(1, CFTimeZoneCreateWithName);
	FUNC_LOAD(1, CFTimeZoneCreateWithTimeIntervalFromGMT);
	FUNC_LOAD(1, CFAbsoluteTimeGetCurrent);
	FUNC_LOAD(1, CFAbsoluteTimeGetGregorianDate);
	FUNC_LOAD(1, CFGregorianDateGetAbsoluteTime);
	FUNC_LOAD(1, CFStringGetTypeID);
	FUNC_LOAD(1, CFDictionaryGetTypeID);
	FUNC_LOAD(1, CFDataGetTypeID);
	FUNC_LOAD(1, CFNumberGetTypeID);
	FUNC_LOAD(1, CFAllocatorGetTypeID);
	FUNC_LOAD(1, CFURLGetTypeID);
	FUNC_LOAD(1, CFReadStreamGetTypeID);
	FUNC_LOAD(1, CFDictionaryReplaceValue);
	FUNC_LOAD(1, CFArrayGetTypeID);
	FUNC_LOAD(1, CFDateGetTypeID);
	FUNC_LOAD(1, CFErrorGetTypeID);
	FUNC_LOAD(1, CFNullGetTypeID);
	FUNC_LOAD(1, CFBooleanGetTypeID);
	FUNC_LOAD(1, CFAttributedStringGetTypeID);
	FUNC_LOAD(1, CFBagGetTypeID);
	FUNC_LOAD(1, CFBitVectorGetTypeID);
	FUNC_LOAD(1, CFBundleGetTypeID);
	FUNC_LOAD(1, CFCalendarGetTypeID);
	FUNC_LOAD(1, CFCharacterSetGetTypeID);
	FUNC_LOAD(1, CFLocaleGetTypeID);
	FUNC_LOAD(1, CFRunArrayGetTypeID);
	FUNC_LOAD(1, CFSetGetTypeID);
	FUNC_LOAD(1, CFTimeZoneGetTypeID);
	FUNC_LOAD(1, CFTreeGetTypeID);
	FUNC_LOAD(1, CFUUIDGetTypeID);
	FUNC_LOAD(1, CFWriteStreamGetTypeID);
	FUNC_LOAD(1, CFXMLNodeGetTypeID);
	FUNC_LOAD(1, CFStorageGetTypeID);
	FUNC_LOAD(1, CFSocketGetTypeID);
	FUNC_LOAD(1, CFWindowsNamedPipeGetTypeID);
	FUNC_LOAD(1, CFPlugInGetTypeID);
	FUNC_LOAD(1, CFPlugInInstanceGetTypeID);
	FUNC_LOAD(1, CFBinaryHeapGetTypeID);
	FUNC_LOAD(1, CFDateFormatterGetTypeID);
	FUNC_LOAD(1, CFMessagePortGetTypeID);
	FUNC_LOAD(1, CFNotificationCenterGetTypeID);
	FUNC_LOAD(1, CFNumberFormatterGetTypeID);
	FUNC_LOAD(1, _CFKeyedArchiverUIDGetTypeID);
	FUNC_LOAD(1, _CFKeyedArchiverUIDGetValue);
	FUNC_LOAD(1, CFStringCreateWithFormat);

#ifdef WIN32
	if (NULL == __CFStringMakeConstantString)
		MessageBoxA(NULL, "CoreFoundation.dll load fun fail", "error", 0);
#endif
}



CCoreFoundation::~CCoreFoundation(void)
{
#ifdef WIN32
	if (m_hDll)
		FreeLibrary(m_hDll);
#endif
}


CMobileDevice::CMobileDevice(bool newdll)
{
#ifdef WIN32
	if (newdll)
		m_hDll = LoadLibrary(L"MobileDevice.dll");
	else
		m_hDll = LoadLibrary(L"iTunesMobileDevice.dll");
	if (NULL == m_hDll)
		MessageBoxA(NULL, "MobileDevice.dll load fail", "error", 0);
#else
	const char* cDylibPath = "/System/Library/PrivateFrameworks/MobileDevice.framework";
	CFStringRef pdllPath = ::CFStringCreateWithCString((::CFAllocatorRef)0, cDylibPath, 134217984);
	if (NULL == pdllPath) 
		return;
	CFURLRef pURLRef = ::CFURLCreateWithFileSystemPath(0, pdllPath, kCFURLPOSIXPathStyle, 0);
	if (NULL == pURLRef) 
		return; 
	CFBundleRef pModle = ::CFBundleCreate(::kCFAllocatorDefault, pURLRef);
	if (NULL == pModle) 
		return; 
#endif
	FUNC_LOAD(0,AMDeviceLookupApplications);
	FUNC_LOAD(0,AMDeviceInstallApplication);
	FUNC_LOAD(0,AMDeviceRemoveApplicationArchive);
	FUNC_LOAD(0,AMDeviceUninstallApplication);
	FUNC_LOAD(0,AMDeviceArchiveApplication);
	FUNC_LOAD(0,AMDeviceStartHouseArrestService);
	FUNC_LOAD(0,AFCConnectionOpen);
	FUNC_LOAD(0,AMDServiceConnectionInvalidate);
	FUNC_LOAD(0,AMDeviceNotificationSubscribe);
	FUNC_LOAD(0,AMDeviceNotificationUnsubscribe);
	FUNC_LOAD(0,AMDeviceRelease);
	FUNC_LOAD(0,AMDeviceConnect);
	FUNC_LOAD(0,AMDeviceDisconnect);
	FUNC_LOAD(0,AMDeviceIsPaired);
	FUNC_LOAD(0,AMDeviceValidatePairing);
	FUNC_LOAD(0,AMDevicePair);
	FUNC_LOAD(0,AMDeviceUnpair);
	FUNC_LOAD(0,AMDeviceStartSession);
	FUNC_LOAD(0,AMDeviceSecureStartService);
	FUNC_LOAD(0,AMDeviceStartService);
	FUNC_LOAD(0,AMDeviceStopSession);
	FUNC_LOAD(0,AFCConnectionClose);
	FUNC_LOAD(0,AFCDeviceInfoOpen);
	FUNC_LOAD(0,AFCFileInfoOpen);
	FUNC_LOAD(0,AFCKeyValueRead);
	FUNC_LOAD(0,AFCKeyValueClose);
	FUNC_LOAD(0,AFCDirectoryOpen);
	FUNC_LOAD(0,AFCDirectoryRead);
	FUNC_LOAD(0,AFCDirectoryClose);
	FUNC_LOAD(0,AFCDirectoryCreate);
	FUNC_LOAD(0,AFCRemovePath);
	FUNC_LOAD(0,AFCRenamePath);
	FUNC_LOAD(0,AFCFileRefOpen);
	FUNC_LOAD(0,AFCFileRefRead);
	FUNC_LOAD(0,AFCFileRefWrite);
	FUNC_LOAD(0,AFCFileRefClose);
	FUNC_LOAD(0,AFCFileRefSeek);
	FUNC_LOAD(0,AFCFileRefTell);
	FUNC_LOAD(0,AMDeviceCopyDeviceIdentifier);
	FUNC_LOAD(0,AMDeviceCopyValue);
	FUNC_LOAD(0,AMDeviceGetInterfaceType);
	FUNC_LOAD(0,AMRestoreRegisterForDeviceNotifications);
	FUNC_LOAD(0,USBMuxConnectByPort);
	FUNC_LOAD(0,AMRestorePerformRecoveryModeRestore);
	FUNC_LOAD(0,AMRestorePerformDFURestore);
	FUNC_LOAD(0,AMRestorableDeviceRegisterForNotificationsForDevices);
	FUNC_LOAD(0,AMRestoreUnregisterForDeviceNotifications);
	FUNC_LOAD(0,AMRestorableDeviceRestore);
	FUNC_LOAD(0,AMSRestoreWithApplications);
	FUNC_LOAD(0,AMSUnregisterTarget);
	FUNC_LOAD(0,AMDeviceSetValue);
	FUNC_LOAD(0,AMRecoveryModeDeviceSendFileToDevice);
	FUNC_LOAD(0,AMRecoveryModeDeviceSendCommandToDevice);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetProductID);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetProductType);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetChipID);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetECID);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetLocationID);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetBoardID);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetProductionMode);
	FUNC_LOAD(0,AMRecoveryModeDeviceGetTypeID);
	FUNC_LOAD(0,AMRecoveryModeGetSoftwareBuildVersion);
	FUNC_LOAD(0,AMDFUModeDeviceGetProductID);
	FUNC_LOAD(0,AMDFUModeDeviceGetProductType);
	FUNC_LOAD(0,AMDFUModeDeviceGetChipID);
	FUNC_LOAD(0,AMDFUModeDeviceGetECID);
	FUNC_LOAD(0,AMDFUModeDeviceGetLocationID);
	FUNC_LOAD(0,AMDFUModeDeviceGetBoardID);
	FUNC_LOAD(0,AMDFUModeDeviceGetProductionMode);
	FUNC_LOAD(0,AMDFUModeDeviceGetTypeID);
	FUNC_LOAD(0,AMRecoveryModeDeviceSetAutoBoot);
	FUNC_LOAD(0,AMRecoveryModeDeviceReboot);
	FUNC_LOAD(0,AMRestoreModeDeviceReboot);
	FUNC_LOAD(0,AMRestoreEnableFileLogging);
	FUNC_LOAD(0,AMRestoreDisableFileLogging);
	FUNC_LOAD(0,AMRestorableDeviceGetState);
	FUNC_LOAD(0,AMRestorableDeviceCopyDFUModeDevice);
	FUNC_LOAD(0,AMRestorableDeviceCopyRecoveryModeDevice);
	FUNC_LOAD(0,AMRestorableDeviceCopyAMDevice);
	FUNC_LOAD(0,AMRestorableDeviceCreateFromAMDevice);
	FUNC_LOAD(0,AMRestorableDeviceGetProductID);
	FUNC_LOAD(0,AMRestorableDeviceGetProductType);
	FUNC_LOAD(0,AMRestorableDeviceGetChipID);
	FUNC_LOAD(0,AMRestorableDeviceGetECID);
	FUNC_LOAD(0,AMRestorableDeviceGetLocationID);
	FUNC_LOAD(0,AMRestorableDeviceGetBoardID);
	FUNC_LOAD(0,AMRestoreModeDeviceGetTypeID);
	FUNC_LOAD(0,AMRestoreModeDeviceCopySerialNumber);
	FUNC_LOAD(0,AMRestorableDeviceCopySerialNumber);
	FUNC_LOAD(0,AMRecoveryModeDeviceCopySerialNumber);
	FUNC_LOAD(0,AFCConnectionGetContext);
	FUNC_LOAD(0,AFCConnectionGetFSBlockSize);
	FUNC_LOAD(0,AFCConnectionGetIOTimeout);
	FUNC_LOAD(0,AFCConnectionGetSocketBlockSize);
	FUNC_LOAD(0,AMRestoreCreateDefaultOptions);
	FUNC_LOAD(0,AMRestorePerformRestoreModeRestore);
	FUNC_LOAD(0,AMRestoreModeDeviceCreate);
	FUNC_LOAD(0,AMRestoreCreatePathsForBundle);
	FUNC_LOAD(0,AMDeviceGetConnectionID);
	FUNC_LOAD(0,AMDeviceEnterRecovery);
	FUNC_LOAD(0,AMDeviceRetain);
	FUNC_LOAD(0,AMDShutdownNotificationProxy);
	FUNC_LOAD(0,AMDeviceDeactivate);
	FUNC_LOAD(0,AMDeviceActivate);
	FUNC_LOAD(0,AMDeviceRemoveValue);
	FUNC_LOAD(0,USBMuxListenerCreate);
	FUNC_LOAD(0,USBMuxListenerHandleData);
	FUNC_LOAD(0,AMDObserveNotification);
	FUNC_LOAD(0,AMSInitialize);
	FUNC_LOAD(0,AMDListenForNotifications);
	FUNC_LOAD(0,AMDeviceStartServiceWithOptions);
	FUNC_LOAD(0,AMDServiceConnectionCreate);
	FUNC_LOAD(0,AMDServiceConnectionGetSocket);
	FUNC_LOAD(0,AMDServiceConnectionGetSecureIOContext);
	FUNC_LOAD(0,AMDServiceConnectionReceive);
	FUNC_LOAD(0,AMDServiceConnectionSend);
	FUNC_LOAD(0,AMDServiceConnectionReceiveMessage);
	FUNC_LOAD(0,AMDServiceConnectionSendMessage);
	FUNC_LOAD(0,AMSChangeBackupPassword);
	FUNC_LOAD(0,AMSBackupWithOptions);
	FUNC_LOAD(0,AMSCancelBackupRestore);
	FUNC_LOAD(0,AMSGetErrorReasonForErrorCode);

#ifdef WIN32
	if (NULL == AMDeviceNotificationSubscribe)
		MessageBoxA(NULL, "MobileDevice.dll load fun fail", "error", 0);
#endif
}

CMobileDevice::~CMobileDevice(void)
{
#ifdef WIN32
	if (m_hDll)
		FreeLibrary(m_hDll);
#endif
}


CAirTrafficHost::CAirTrafficHost(void)
{
#ifdef WIN32
	m_hDll = LoadLibrary(L"AirTrafficHost.dll");
	if (NULL == m_hDll)
		MessageBoxA(NULL, "AirTrafficHost.dll load fail", "error", 0);
#else
	const char* cDylibPath = "/System/Library/PrivateFrameworks/AirTrafficHost.framework";
	CFStringRef pdllPath = ::CFStringCreateWithCString((::CFAllocatorRef)0, cDylibPath, 134217984);
	if (NULL == pdllPath)
		return;
	CFURLRef pURLRef = ::CFURLCreateWithFileSystemPath(0, pdllPath, kCFURLPOSIXPathStyle, 0);
	if (NULL == pURLRef)
		return;
	CFBundleRef pModle = ::CFBundleCreate(::kCFAllocatorDefault, pURLRef);
	if (NULL == pModle)
		return;
#endif
	FUNC_LOAD(0, ATCFMessageGetParam);
	FUNC_LOAD(0, ATHostConnectionCreate);
	FUNC_LOAD(0, ATHostConnectionCreateWithLibrary);
	FUNC_LOAD(0, ATHostConnectionSendPing);
	FUNC_LOAD(0, ATHostConnectionSendAssetMetricsRequest);
	FUNC_LOAD(0, ATHostConnectionInvalidate);
	FUNC_LOAD(0, ATHostConnectionClose);
	FUNC_LOAD(0, ATHostConnectionGetCurrentSessionNumber);
	FUNC_LOAD(0, ATHostConnectionRelease);
	FUNC_LOAD(0, ATHostConnectionSendPowerAssertion);
	FUNC_LOAD(0, ATHostConnectionRetain);
	FUNC_LOAD(0, ATHostConnectionSendMetadataSyncFinished);
	FUNC_LOAD(0, ATHostConnectionSendFileError);
	FUNC_LOAD(0, ATCFMessageCreate);
	FUNC_LOAD(0, ATHostConnectionSendFileProgress);
	FUNC_LOAD(0, ATHostConnectionSendAssetCompleted);
	FUNC_LOAD(0, ATCFMessageGetName);
	FUNC_LOAD(0, ATHostConnectionSendHostInfo);
	FUNC_LOAD(0, ATHostConnectionSendSyncRequest);
	FUNC_LOAD(0, ATHostConnectionSendMessage);
	FUNC_LOAD(0, ATHostConnectionGetGrappaSessionId);
	FUNC_LOAD(0, ATHostConnectionReadMessage);
#ifdef WIN32
	if (NULL == ATHostConnectionReadMessage)
		MessageBoxA(NULL, "AirTrafficHost.dll load fun fail", "error", 0);
#endif
}

CAirTrafficHost::~CAirTrafficHost(void)
{
#ifdef WIN32
	if (m_hDll)
		FreeLibrary(m_hDll);
#endif
}



