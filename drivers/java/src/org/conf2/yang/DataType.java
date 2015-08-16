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
}
