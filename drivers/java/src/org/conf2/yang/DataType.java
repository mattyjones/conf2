package org.conf2.yang;

/**
 *
 */
public class DataType extends MetaBase {
    public String range;
    public int minLength;
    public int maxLength;
    public String path;
    public String pattern;
    public String[] enumeration;
    public ValueType valType;

    public DataType(String ident) {
        super(ident);
    }

    @Override
    public void setIdent(String ident) {
        if ("int32".equals(ident)) {
            valType = ValueType.INT32;
        } else if ("string".equals(ident)) {
            valType = ValueType.STRING;
        } else if ("boolean".equals(ident)) {
            valType = ValueType.BOOLEAN;
        }
    }
}
