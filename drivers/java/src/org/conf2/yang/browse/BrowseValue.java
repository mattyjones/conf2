package org.conf2.yang.browse;

import org.conf2.yang.ValueType;

import java.util.Arrays;

/**
 *
 */
public class BrowseValue {
    public ValueType valType = ValueType.EMPTY;
    public boolean isList;
    public String str;
    public boolean bool;
    public int int32;
    public String[] strlist;
    public boolean[] boollist;
    public int[] int32list;

    public static BrowseValue Strlist(String[] v) {
        BrowseValue bv = new BrowseValue();
        bv.strlist = v;
        bv.valType = ValueType.STRING;
        bv.isList = true;
        return bv;
    }

    public static BrowseValue Int32list(int[] v) {
        BrowseValue bv = new BrowseValue();
        bv.int32list = v;
        bv.valType = ValueType.INT32;
        bv.isList = true;
        return bv;
    }

    public static BrowseValue Boollist(boolean[] v) {
        BrowseValue bv = new BrowseValue();
        bv.boollist = v;
        bv.valType = ValueType.BOOLEAN;
        bv.isList = true;
        return bv;
    }

    public static BrowseValue Str(String v) {
        BrowseValue bv = new BrowseValue();
        bv.str = v;
        bv.valType = ValueType.STRING;
        return bv;
    }

    public static BrowseValue Int32(int v) {
        BrowseValue bv = new BrowseValue();
        bv.int32 = v;
        bv.valType = ValueType.INT32;
        return bv;
    }

    public static BrowseValue Bool(boolean v) {
        BrowseValue bv = new BrowseValue();
        bv.bool = v;
        bv.valType = ValueType.BOOLEAN;
        return bv;
    }

    public static BrowseValue fromByteArray(int len, int typeCode, byte[] data) {
System.out.printf("fromByteArray\n");
        ValueType type = ValueType.values()[typeCode];
        BrowseValue v = new BrowseValue();
        v.valType = type;
        v.isList = true;
        switch (type) {
            case STRING:
System.out.printf("string array len=%d, buflen=%d\n", len, data.length);
                strlistFromByteArray(v, len, data);
                break;
            case INT32:
                intlistFromByteArray(v, len, data);
                break;
            case BOOLEAN:
                booleanlistFromByteArray(v, len, data);
                break;
        }

        return v;
    }

    static void strlistFromByteArray(BrowseValue v, int len, byte[] data) {
        v.strlist = new String[len];
        int substrStart = 0;
        for (int i = 0; i < len; i++) {
            int substrEnd = substrStart + 1;
            while (substrEnd < data.length && data[substrEnd] != '\0') {
System.out.printf("strlistFromByteArray, %d, %c\n", substrEnd, data[substrEnd]);
                substrEnd += 1;
            }
System.out.printf("BEFORE string substrStart=%d, substrEnd=%d\n", substrStart, substrEnd);
            v.strlist[i] = new String(data, substrStart, substrEnd);
System.out.printf("AFTER string substrStart=%d, substrEnd=%d\n", substrStart, substrEnd);
        }
    }

    static void intlistFromByteArray(BrowseValue v, int len, byte[] data) {
        v.int32list = new int[len];
        for (int i = 0; i < len; i += 4) {
            // big endian. TODO: little-endian
            v.int32list[i] = ((int)data[i]) << 24;
            v.int32list[i] += ((int)data[i + 1]) << 16;
            v.int32list[i] += ((int)data[i + 2]) << 8;
            v.int32list[i] += ((int)data[i + 3]);
        }
    }

    static void booleanlistFromByteArray(BrowseValue v, int len, byte[] data) {
        v.boollist = new boolean[len];
        for (int i = 0; i < len; i += 2) {
            v.boollist[i] = (data[i] > 0);
            v.boollist[i] = v.boollist[i] || (data[i + 1] > 0);
        }
    }

    int decodeValueType() {
        return valType.code;
    }
}
