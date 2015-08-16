package org.conf2.yang.browse;

import org.conf2.yang.DataType;
import org.conf2.yang.ValueType;

/**
 *
 */
public class BrowseValue {
    public ValueType valType = ValueType.EMPTY;
    public String str;
    public boolean bool;
    public int int32;

    public static BrowseValue Str(String v) {
        BrowseValue bv = new BrowseValue();
        bv.str = v;
        bv.valType = ValueType.STRING;
        return bv;
    }

    public static BrowseValue Int32(String v) {
        BrowseValue bv = new BrowseValue();
        bv.str = v;
        bv.valType = ValueType.INT32;
        return bv;
    }

    public static BrowseValue Bool(boolean v) {
        BrowseValue bv = new BrowseValue();
        bv.bool = v;
        bv.valType = ValueType.BOOLEAN;
        return bv;
    }
}
