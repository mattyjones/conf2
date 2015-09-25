package org.conf2.schema.comm;

import org.conf2.schema.browse.BrowseValue;

import java.io.IOException;
import java.nio.ByteBuffer;

/**
 *
 */
public class DataWriter {
    static byte FALSE_BYTE = 0;
    static byte TRUE_BYTE = 1;
    private ByteBuffer data;

    public void writeValue(BrowseValue val) {

        // finding the size and allocating buffer once allows C code to access byte array
        // directly and pass along to Go to directly read.
        int encodedValueSize = determineEncodedValueSize(val);
        data = ByteBuffer.allocateDirect(encodedValueSize);

        switch (val.valType) {
            case ENUMERATION:
            case INT32:
                data.putInt(val.int32);
                break;
            case BOOLEAN:
                data.put(val.bool ? TRUE_BYTE : FALSE_BYTE);
                break;
            case STRING:
                data.put(val.str.getBytes());
                data.put(DataReader.CSTR_TERM);
                break;
            case ENUMERATION_LIST:
            case INT32_LIST:
                data.putInt(val.listLen());
                for ( int i : val.int32list) {
                    data.putInt(i);
                }
                break;
            case BOOLEAN_LIST:
                data.putInt(val.listLen());
                for (boolean b : val.boollist) {
                    data.putInt(b ? TRUE_BYTE : FALSE_BYTE);
                }
                break;
            case STRING_LIST:
                data.putInt(val.listLen());
                for (String s : val.strlist) {
                    data.put(s.getBytes());
                    data.put(DataReader.CSTR_TERM);
                }
                break;
        }
    }

    public int determineEncodedValueSize(BrowseValue val) {
        int nullTerminatorSize = 1;
        int intSize = 4;
        int formatCodeSize = intSize; // format code
        int listLenSize = intSize;
        int size = formatCodeSize;
        switch (val.valType) {
            case ENUMERATION:
            case INT32:
                size += 4;
                break;
            case BOOLEAN:
                size += 1;
                break;
            case STRING:
                size += utf8StrLen(val.str) + nullTerminatorSize;
                break;
            case ENUMERATION_LIST:
            case INT32_LIST:
                size += listLenSize + (val.listLen() * intSize);
                break;
            case BOOLEAN_LIST:
                size += listLenSize + val.listLen();
                break;
            case STRING_LIST:
                size += listLenSize;
                for (String s : val.strlist) {
                    size += utf8StrLen(s) + nullTerminatorSize;
                }
                break;
        }
        return size;
    }

    public ByteBuffer getByteBuffer() {
        return data;
    }

    public void writeString(String s) throws IOException {
        data.put(s.getBytes());
        data.put(DataReader.CSTR_TERM);
    }

    /**
     * Does not include null termination byte
     */
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
}
