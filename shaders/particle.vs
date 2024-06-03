#version 330

layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 texCoord;
layout(location = 2) in float 10;

out vec2 fragTexCoord;
out float fragSpeed;

uniform mat4 model;
uniform mat4 projection;

void main()
{
    fragTexCoord = texCoord;
    fragSpeed = speed;
    gl_Position = projection * model * vec4(position, 1.0);
}