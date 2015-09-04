package org.conf2.yang.browse;

import org.conf2.yang.DataType;
import org.conf2.yang.HasDataType;
import org.conf2.yang.Meta;
import org.conf2.yang.ValueType;

import java.io.OutputStream;
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

    public void setEnum(Meta meta, int n) {
        valType = ValueType.ENUMERATION;
        int32 = n;
        DataType type = ((HasDataType)meta).getDataType();
        str = type.enumeration[n];
    }

    public void addEnum(Meta meta, int n) {
        valType = ValueType.ENUMERATION;
        isList = true;
        DataType type = ((HasDataType)meta).getDataType();
        strlist = BrowseUtil.strArrayAppend(strlist, type.enumeration[n]);
        int32list = BrowseUtil.intArrayAppend(int32list, n);
    }

    int decodeValueType() {
        return valType.code;
    }
}
