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
    UNION(18),

    BINARY_LIST(1025),
    BITS_LIST(1026),
    BOOLEAN_LIST(1027),
    DECIMAL64_LIST(1028),
    ENUMERATION_LIST(1029),
    IDENTITYDEF_LIST(1030),
    INSTANCE_IDENTIFIER_LIST(1031),
    INT8_LIST(1032),
    INT16_LIST(1033),
    INT32_LIST(1034),
    INT64_LIST(1035),
    LEAFREF_LIST(1036),
    STRING_LIST(1037),
    UINT8_LIST(1038),
    UINT16_LIST(1039),
    UINT32_LIST(1040),
    UINT64_LIST(1041),
    UNION_LIST(1042);

    public final int code;

    ValueType(int code) {
        this.code = code;
    }
}
