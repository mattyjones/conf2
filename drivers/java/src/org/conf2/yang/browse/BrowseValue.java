package org.conf2.yang.browse;

import org.conf2.yang.ValueType;

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
        ValueType type = ValueType.values()[typeCode];
        BrowseValue v = new BrowseValue();
        v.valType = type;
        v.isList = true;
        switch (type) {
            case STRING:
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
            while (data[substrEnd] != '\0') {
                substrEnd += 1;
            }
            v.strlist[i] = new String(data, substrStart, substrEnd);
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
}
