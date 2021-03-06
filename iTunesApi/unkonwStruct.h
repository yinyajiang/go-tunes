#pragma once
#ifdef __cplusplus
extern "C"
{
#endif
#ifndef WIN32
#include <CoreFoundation/CoreFoundation.h>
#endif
#include <stdint.h>
#pragma pack(push, 1)

/*Messages passed to device notification callbacks */
#define ADNCI_MSG_CONNECTECD 1
#define ADNCI_MSG_DISCONNECTED 2
#define ADNCI_MSG_UNKNOWN 3
#define AMD_IPHONE_PRODUCT_ID 0x1290

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
#ifdef WIN32
	typedef void *CFStringRef;
	typedef void *CFTypeRef;
	typedef void *CFMutableDataRef;
	typedef void *CFPropertyListRef;
	typedef void *CFDictionaryRef;
	typedef void *CFMutableDictionaryRef;
	typedef void *CFDataRef;
	typedef void *CFNumberRef;
	typedef void *CFAllocatorRef;
	typedef void *CFURLRef;
	typedef void *CFReadStreamRef;
	typedef void *CFArrayRef;
	typedef void *CFDateRef;
	typedef void *CFErrorRef;
	typedef void *CFMutableArrayRef;
	typedef void *CFBooleanRef;
	typedef void *CFTimeZoneRef;
	typedef int CFIndex;
	typedef int CFTypeID;
	typedef int CFStringEncoding;
	typedef int CFURLPathStyle; // posix=0, hfs=1, windows=2
	typedef void *CFBundleRef;
	typedef void *CFRunLoopRef;
	typedef double CFAbsoluteTime;
	typedef double CFTimeInterval;
#ifdef WIN32
	typedef unsigned char Boolean;
	typedef uint8_t UInt8;
#endif

	typedef struct _CFRange
	{
		CFIndex location;
		CFIndex length;
	} CFRange;

	typedef struct _CFGregorianDate
	{
		int32_t year;
		int8_t month;
		int8_t day;
		int8_t hour;
		int8_t minute;
		double second;
	} CFGregorianDate;

	typedef struct _CFGregorianUnits
	{
		int32_t years;
		int32_t months;
		int32_t days;
		int32_t hours;
		int32_t minutes;
		double seconds;
	} CFGregorianUnits;

	enum
	{
		kCFStringEncodingInvalidId = -1L,
		kCFStringEncodingMacRoman = 0L,
		kCFStringEncodingWindowsLatin1 = 0x0500, /* ANSI codepage 1252 */
		kCFStringEncodingISOLatin1 = 0x0201,	 /* ISO 8850 1 */
		kCFStringEncodingNextStepLatin = 0x0B01, /* NextStep encoding*/
		kCFStringEncodingASCII = 0x0600,		 /* 0..127 */
		kCFStringEncodingUnicode = 0x0100,		 /* kTextEncodingUnicodeDefault  + kTextEncodingDefaultFormat (aka kUnicode16BitFormat) */
		kCFStringEncodingUTF8 = 0x08000100,		 /* kTextEncodingUnicodeDefault + kUnicodeUTF8Format */
		kCFStringEncodingNonLossyASCII = 0x0BFF	 /* 7bit Unicode variants used by YellowBox & Java */
	};

	enum
	{
		kCFPropertyListImmutable = 0,
		kCFPropertyListMutableContainers,
		kCFPropertyListMutableContainersAndLeaves
	};

	typedef enum _CFPropertyListFormat
	{
		kCFPropertyListOpenStepFormat = 1,
		kCFPropertyListXMLFormat_v1_0 = 100,
		kCFPropertyListBinaryFormat_v1_0 = 200
	} CFPropertyListFormat;

	typedef enum _CFNumberType
	{
		/* Types from MacTypes.h */
		kCFNumberSInt8Type = 1,
		kCFNumberSInt16Type = 2,
		kCFNumberSInt32Type = 3,
		kCFNumberSInt64Type = 4,
		kCFNumberFloat32Type = 5,
		kCFNumberFloat64Type = 6, /* 64-bit IEEE 754 */
		/* Basic C types */
		kCFNumberCharType = 7,
		kCFNumberShortType = 8,
		kCFNumberIntType = 9,
		kCFNumberLongType = 10,
		kCFNumberLongLongType = 11,
		kCFNumberFloatType = 12,
		kCFNumberDoubleType = 13, /* Other */
		kCFNumberCFIndexType = 14,
		kCFNumberNSIntegerType = 15,
		kCFNumberMaxType = 16
	} CFNumberType;
#endif

	struct am_recovery_device;
	struct am_dfu_device;
	struct am_restore_device;
	struct am_device_notification_callback_info;

	typedef void (*BackupCallBack)(CFStringRef targetID, int percentOrErrorCode, int cookie, CFStringRef backupPath, void *, void *, void *, void *, void *, void *, void *, void *);
	typedef void (*am_restore_device_notification_callback)(struct am_recovery_device *, void *);
	typedef void (*am_dfu_device_notification_callback)(struct am_dfu_device *, void *);
	typedef void (*am_recovery_device_notification_callback)(struct am_restore_device *, int);
	typedef void (*am_device_notification_callback)(struct am_device_notification_callback_info *);

	struct am_recovery_device
	{
		unsigned char unk[8];
		am_restore_device_notification_callback callback;
		void *user_info;
		unsigned char unk1[12];
		unsigned int readwrite_pipe;
		unsigned char read_pipe;
		unsigned char write_ctrl_pipe;
		unsigned char read_unk_pipe;
		unsigned char write_unk_pipe;
		unsigned char write_input_pipe;
	};

	struct am_restore_device
	{
		unsigned char unk[32];
		int port;
	};

	struct am_dfu_device
	{
		unsigned char unk[32];
		int port;
	};

	struct am_device
	{
		unsigned char unk0[16];
		unsigned int device_id;
		unsigned int product_id;
		char *serial;
		unsigned int unk1;
		unsigned char unk2[4];
		unsigned int lockdown_conn;
		unsigned char unk3[8];
	};

	struct am_device_notification
	{
		unsigned int unk0;
		unsigned int unk1;
		unsigned int unk2;
		am_device_notification_callback callback;
		unsigned int unk3;
	};

#pragma pack(pop)

#ifdef __cplusplus
};
#endif
