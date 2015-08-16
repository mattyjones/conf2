package org.conf2.yang.driver;

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
