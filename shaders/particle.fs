#version 330

in vec2 fragTexCoord;
in float fragSpeed;

out vec4 finalColor;

uniform sampler2D texture0;

void main()
{
    vec4 texelColor = texture(texture0, fragTexCoord);

    // Calculate color based on speed
    float intensity = clamp(fragSpeed / 10.0, 0.0, 1.0);
    vec3 color = mix(vec3(0.0, 0.0, 1.0), vec3(1.0, 0.0, 0.0), intensity); // Blue to Red gradient

    finalColor = vec4(color, texelColor.a);
}