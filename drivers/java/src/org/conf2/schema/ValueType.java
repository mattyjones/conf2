package org.conf2.schema;

/**
 *
 */
public enum ValueType {
    // In specific order - see browse.h
    EMPTY(0),
    BINARY(1),
    BITS(2),
    BOOLEAN(3),
    DECIMAL64(4),
    ENUMERATION(5),
    IDENTITYDEF(6),
    INSTANCE_IDENTIFIER(7),
    INT8(8),
    INT16(9),
    INT32(10),
    INT64(11),
    LEAFREF(12),
    STRING(13),
    UINT8(14),
    UINT16(15),
    UINT32(16),
    UINT64(17),
    UNION(18);

    public final int code;

    ValueType(int code) {
        this.code = code;
    }
}
