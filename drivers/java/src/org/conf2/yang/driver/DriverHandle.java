package org.conf2.yang.driver;

import java.nio.ByteBuffer;
import java.util.Arrays;

/**
 *
 */
public class DriverHandle {
    public ByteBuffer reference;
    public DriverHandle(byte[] reference) {
        this.reference = ByteBuffer.allocateDirect(reference.length);
        this.reference.put(reference);
    }
}
