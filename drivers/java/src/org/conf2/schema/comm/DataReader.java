package org.conf2.schema.comm;

import org.conf2.schema.DataType;
import org.conf2.schema.ValueType;
import org.conf2.schema.browse.BrowseValue;

import java.io.IOException;
import java.nio.ByteBuffer;

/**
 *
 */
public class DataReader {
    static final byte CSTR_TERM = 0;
    private ByteBuffer data;

    public DataReader(ByteBuffer data) {
        this.data = data;
    }

    public BrowseValue[] readValues(DataType[] expectedDataTypes)  {
        int valueLen = data.getInt();
        if (valueLen != expectedDataTypes.length) {
            throw new RuntimeException(String.format("Expected %d values but only passed %d",
                    expectedDataTypes.length, valueLen));
        }
        BrowseValue[] vals = new BrowseValue[valueLen];
        for (int i = 0; i < valueLen; i++) {
            vals[i] = readValue(expectedDataTypes[i]);
        }
        return null;
    }

    public BrowseValue readValue(DataType expectedDataType) {
        int formatCode = data.getInt();
        ValueType actualDataType = ValueType.values()[formatCode];
        if (actualDataType != expectedDataType.valType) {
            throw new RuntimeException(String.format("Unexpected value type %s. Expected %s",
                    actualDataType.name(), expectedDataType.valType.name()));
        }
        BrowseValue v = new BrowseValue();
        int listLen = 0;
        switch (actualDataType) {
            case INT32:
                v.int32 = data.getInt();
                break;
            case ENUMERATION:
                int enumId = data.getInt();
                v.setEnum(expectedDataType, enumId);
                break;
            case BOOLEAN:
                v.bool = data.get() == DataWriter.TRUE_BYTE;
            case STRING:
                v.str = readString();
                break;
            case ENUMERATION_LIST:
            case INT32_LIST:
                listLen = data.getInt();
                v.int32list = new int[listLen];
                for (int i = 0; i < listLen; i++) {
                    v.int32list[i] = data.getInt();
                }
                if (actualDataType == ValueType.ENUMERATION_LIST) {
                    v.setEnumList(expectedDataType, v.int32list);
                }
                break;
            case BOOLEAN_LIST:
                listLen = data.getInt();
                v.boollist = new boolean[listLen];
                for (int i = 0; i < listLen; i++) {
                    v.boollist[i] = data.get() == DataWriter.TRUE_BYTE;
                }
                break;
            case STRING_LIST:
                listLen = data.getInt();
                v.strlist = new String[listLen];
                for (int i = 0; i < listLen; i++) {
                    v.strlist[i] = readString();
                }
                break;

        }
        return null;
    }

    public String readString() {
        int strStart = data.position();
        int strEnd = strStart;
        for (;strEnd < data.capacity() && data.get() != 0; strEnd++);
        int strLen = (data.position() - 1) - strStart;
        byte[] bytes = new byte[strLen];
        data.position(strStart);
        data.get(bytes);
        String s = new String(bytes);
        data.get(); // null term
        return s;
    }
}
