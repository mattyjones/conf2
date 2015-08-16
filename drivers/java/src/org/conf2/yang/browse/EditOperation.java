package org.conf2.yang.browse;

/**
 *
 */
public enum EditOperation {
    UNKNOWN,
    CREATE_CHILD,
    POST_CREATE_CHILD,
    CREATE_LIST,                       // 3
    POST_CREATE_LIST,                  // 4
    UPDATE_VALUE,                      // 5
    DELETE_CHILD,                      // 6
    DELETE_LIST,                       // 7
    BEGIN_EDIT,                        // 8
    END_EDIT,                          // 9
    CREATE_LIST_ITEM,                  // 10
    POST_CREATE_LIST_ITEM,             // 11
}
