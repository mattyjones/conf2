package org.conf2.schema.browse;

import org.conf2.schema.DataType;
import org.conf2.schema.HasDataType;
import org.conf2.schema.Meta;
import org.conf2.schema.ValueType;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.ByteBuffer;
import java.util.Arrays;

/**
 *
 */
public class BrowseValue {
    public ValueType valType = ValueType.EMPTY;
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

    public void setEnum(DataType type, int n) {
        valType = ValueType.ENUMERATION;
        int32 = n;
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

    public void setEnumList(DataType type, int[] enumIds) {
        int32list = enumIds;
        strlist = new String[enumIds.length];
        for (int i = 0; i < enumIds.length; i++) {
            strlist[i] = type.enumeration[int32list[i]];
        }
    }

    public void addEnum(DataType type, int n) {
        valType = ValueType.ENUMERATION_LIST;
        strlist = BrowseUtil.strArrayAppend(strlist, type.enumeration[n]);
        int32list = BrowseUtil.intArrayAppend(int32list, n);
    }

    public int listLen() {
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
}
