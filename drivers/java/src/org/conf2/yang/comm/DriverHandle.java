package org.conf2.yang.comm;

import java.util.Arrays;

/**
 *
 */
public class DriverHandle {
    public byte[] reference;
    public DriverHandle(byte[] reference) {
        this.reference = Arrays.copyOf(reference, reference.length);
    }
}
