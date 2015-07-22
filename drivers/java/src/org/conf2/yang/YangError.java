package org.conf2.yang;

/**
 *
 */
public class YangError extends RuntimeException {
    public YangError(String msg) {
        super(msg);
    }
    public YangError(String msg, Throwable suberr) {
        super(msg, suberr);
    }
}
