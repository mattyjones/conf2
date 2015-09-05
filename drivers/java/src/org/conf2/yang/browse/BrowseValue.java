package org.conf2.yang.browse;

import org.conf2.yang.DataType;
import org.conf2.yang.HasDataType;
import org.conf2.yang.Meta;
import org.conf2.yang.ValueType;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.ByteBuffer;
import java.util.Arrays;

/**
 *
 */
public class BrowseValue {
    private static final byte CSTR_TERM = 0;
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

    public static BrowseValue Enum(int n, String s) {
        BrowseValue bv = new BrowseValue();
        bv.int32 = n;
        bv.str = s;
        bv.valType = ValueType.ENUMERATION;
        return bv;
    }

    public static BrowseValue decodeList(ByteBuffer data, int listLen, int valTypeCode) {
        BrowseValue val = new BrowseValue();
        val.valType = ValueType.values()[valTypeCode];
        val.isList = true;
        switch (val.valType) {
            case ENUMERATION:
                val.int32list = new int[listLen];

                // TODO : String vals

                break;
            case STRING:
                val.strlist = decodeCStrArray(data, listLen);
                break;
            case INT32:
                val.int32list = decodeCIntArray(data, listLen);
                break;
            case BOOLEAN:
                val.boollist = decodeCBoolArray(data, listLen);
                break;
        }
        return val;
    }

    public void setEnum(Meta meta, int n) {
        valType = ValueType.ENUMERATION;
        int32 = n;
        DataType type = ((HasDataType)meta).getDataType();
        str = type.enumeration[n];
    }

    public void setEnum(Meta meta, String enumLabel) {
        valType = ValueType.ENUMERATION;
        DataType type = ((HasDataType)meta).getDataType();
        for (int n= 0; n < type.enumeration.length; n++) {
            if (type.enumeration[n].equals(enumLabel)) {
                int32 = n;
                str = type.enumeration[n];
                return;
            }
        }
    }

    public void addEnum(Meta meta, int n) {
        valType = ValueType.ENUMERATION;
        isList = true;
        DataType type = ((HasDataType)meta).getDataType();
        strlist = BrowseUtil.strArrayAppend(strlist, type.enumeration[n]);
        int32list = BrowseUtil.intArrayAppend(int32list, n);
    }

    public static int utf8StrLen(CharSequence s) {
        int count = 0;
        for (int i = 0, len = s.length(); i < len; i++) {
            char ch = s.charAt(i);
            if (ch <= 0x7F) {
                count++;
            } else if (ch <= 0x7FF) {
                count += 2;
            } else if (Character.isHighSurrogate(ch)) {
                count += 4;
                ++i;
            } else {
                count += 3;
            }
        }
        return count;
    }

    public ByteBuffer encodeList() throws IOException {
        switch (valType) {
            case STRING:
                if (strlist != null) {
                    return encodeCStrArray(strlist);
                }
                break;
            case ENUMERATION:
            case INT32:
                if (int32list != null) {
                    return encodeCIntArray(int32list);
                }
                break;
            case BOOLEAN:
                if (boollist != null) {
                    return encodeCBoolArray(boollist);
                }
                break;
        }
        return null;
    }

    public static String[] decodeCStrArray(ByteBuffer buff, int listlen) {
        String[] strlist = new String[listlen];
        int strStart = 0;
        int strEnd = 0;
        for (int i = 0; strEnd < buff.capacity() && i < listlen; i++) {
            strStart = buff.position();
            for (;strEnd < buff.capacity() && buff.get() != 0; strEnd++) {
            }
            int strLen = (buff.position() - 1) - strStart;
            byte[] bytes = new byte[strLen];
            buff.position(strStart);
            buff.get(bytes);
            strlist[i] = new String(bytes);
            buff.get(); // null term
        }
        return strlist;
    }

    public static boolean[] decodeCBoolArray(ByteBuffer buff, int listlen) {
        boolean[] boollist = new boolean[listlen];
        for (int i = 0; i < listlen; i++) {
            boollist[i] =  buff.getShort() > 0;
        }
        return boollist;
    }

    public static int[] decodeCIntArray(ByteBuffer buff, int listlen) {
        int[] intlist = new int[listlen];
        for (int i = 0; i < listlen; i++) {
            intlist[i] =  buff.getInt();
        }
        return intlist;
    }

    public static ByteBuffer encodeCStrArray(String[] strlist) throws IOException {
        int datalen = 0;
        for (String s : strlist) {
            datalen += utf8StrLen(s) + 1;
        }
        ByteBuffer out = ByteBuffer.allocateDirect(datalen);
        for (String s : strlist) {
            out.put(s.getBytes());
            out.put(CSTR_TERM);
        }
        return out;
    }

    public static ByteBuffer encodeCIntArray(int[] intlist) throws IOException {
        ByteBuffer out = ByteBuffer.allocateDirect(4 * intlist.length);
        for (int i = 0; i < intlist.length; i++) {
            out.putInt(intlist[i]);
        }
        return out;
    }

    public static ByteBuffer encodeCBoolArray(boolean[] boollist) throws IOException {
        short sTrue = 1;
        short sFalse = 0;
        ByteBuffer out = ByteBuffer.allocateDirect(2 * boollist.length);
        for (int i = 0; i < boollist.length; i++) {
            out.putShort(boollist[i] ? sTrue : sFalse);
        }
        return out;
    }

    int listLen() {
        switch (valType) {
            case ENUMERATION:
            case INT32:
                return int32list.length;
            case STRING:
                return strlist.length;
            case BOOLEAN:
                return boollist.length;
        }
        return 0;
    }

    int decodeValueType() {
        return valType.code;
    }
}
